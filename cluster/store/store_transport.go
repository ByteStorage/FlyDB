package store

import (
	"context"
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/encoding"
	raftPB "github.com/ByteStorage/FlyDB/lib/proto/raft"
	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
	"io"
	"sync"
	"time"
)

type ClientConn struct {
	conn   *grpc.ClientConn
	client raftPB.RaftServiceClient
	mtx    sync.Mutex
}
type raftPipeline struct {
	stream        raftPB.RaftService_AppendEntriesPipelineClient
	cancel        func()
	inflightChMtx sync.Mutex
	inflightCh    chan *appendFuture
	doneCh        chan raft.AppendFuture
}

type appendFuture struct {
	raft.AppendFuture
	start    time.Time
	request  *raft.AppendEntriesRequest
	response raft.AppendEntriesResponse
	err      error
	done     chan struct{}
}

// Transport implements raft.Transport interface
// we can use it to send rpc to other raft nodes
// and receive rpc from other raft nodes
type Transport struct {
	//implement me
	localAddr        raft.ServerAddress
	consumer         chan raft.RPC
	clients          map[raft.ServerAddress]*ClientConn
	server           *grpc.Server
	heartbeatFn      func(raft.RPC)
	dialOptions      []grpc.DialOption
	heartbeatTimeout time.Duration
	sync.RWMutex
}

// NewTransport returns a new transport, it needs start a grpc server
func newTransport(conf config.Config) raft.Transport {
	return &Transport{
		localAddr:        conf.LocalAddress,
		dialOptions:      []grpc.DialOption{grpc.WithInsecure()},
		heartbeatTimeout: conf.HeartbeatTimeout,
		consumer:         make(chan raft.RPC),
		clients:          map[raft.ServerAddress]*ClientConn{},
	}
}

// AppendEntriesPipeline returns an interface that can be used to pipeline
// AppendEntries requests.
func (t *Transport) AppendEntriesPipeline(id raft.ServerID, target raft.ServerAddress) (raft.AppendPipeline, error) {
	c, err := t.getPeer(target)
	if err != nil {
		return nil, err
	}
	ctx := context.TODO()
	ctx, cancel := context.WithCancel(ctx)
	stream, err := c.AppendEntriesPipeline(ctx)
	if err != nil {
		cancel()
		return nil, err
	}
	rpa := raftPipeline{
		stream:     stream,
		cancel:     cancel,
		inflightCh: make(chan *appendFuture, 20),
		doneCh:     make(chan raft.AppendFuture, 20),
	}
	go rpa.receiver()
	return &rpa, nil
}

// AppendEntries sends the appropriate RPC to the target node.
func (t *Transport) AppendEntries(id raft.ServerID, target raft.ServerAddress, args *raft.AppendEntriesRequest, resp *raft.AppendEntriesResponse) error {
	c, err := t.getPeer(target)
	if err != nil {
		return err
	}
	ctx := context.TODO()
	if t.heartbeatTimeout > 0 && isHeartbeat(args) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, t.heartbeatTimeout)
		defer cancel()
	}
	ret, err := c.AppendEntries(ctx, encoding.EncodeAppendEntriesRequest(args))
	if err != nil {
		return err
	}
	*resp = *encoding.DecodeAppendEntriesResponse(ret)
	return nil
}

// RequestVote sends the appropriate RPC to the target node.
func (t *Transport) RequestVote(id raft.ServerID, target raft.ServerAddress, args *raft.RequestVoteRequest, resp *raft.RequestVoteResponse) error {
	c, err := t.getPeer(target)
	if err != nil {
		return err
	}
	vote, err := c.RequestVote(context.TODO(), encoding.EncodeRequestVoteRequest(args))
	if err != nil {
		return err
	}
	*resp = *encoding.DecodeRequestVoteResponse(vote)
	return nil
}

// InstallSnapshot is used to push a snapshot down to a follower. The data is read from
// the ReadCloser and streamed to the client.
func (t *Transport) InstallSnapshot(id raft.ServerID, target raft.ServerAddress, args *raft.InstallSnapshotRequest, resp *raft.InstallSnapshotResponse, data io.Reader) error {
	c, err := t.getPeer(target)
	if err != nil {
		return err
	}
	inSnap, err := c.InstallSnapshot(context.TODO(), encoding.EncodeInstallSnapshotRequest(args))
	if err != nil {
		return err
	}

	*resp = *encoding.DecodeInstallSnapshotResponse(inSnap)
	return nil
}

// TimeoutNow is used to start a leadership transfer to the target node.
func (t *Transport) TimeoutNow(id raft.ServerID, target raft.ServerAddress, args *raft.TimeoutNowRequest, resp *raft.TimeoutNowResponse) error {
	c, err := t.getPeer(target)
	if err != nil {
		return err
	}
	ret, err := c.TimeoutNow(context.TODO(), encoding.EncodeTimeoutNowRequest(args))
	if err != nil {
		return err
	}
	*resp = *encoding.DecodeTimeoutNowResponse(ret)
	return nil
}

// Consumer returns a channel that can be used to
// consume and respond to RPC requests.
func (t *Transport) Consumer() <-chan raft.RPC {
	return t.consumer
}

// LocalAddr is used to return our local address to distinguish from our peers.
func (t *Transport) LocalAddr() raft.ServerAddress {
	return t.localAddr
}

// EncodePeer is used to serialize a peer's address.
func (t *Transport) EncodePeer(id raft.ServerID, addr raft.ServerAddress) []byte {
	return []byte(addr)
}

// DecodePeer is used to deserialize a peer's address.
func (t *Transport) DecodePeer(p []byte) raft.ServerAddress {
	return raft.ServerAddress(p)
}

// SetHeartbeatHandler is used to setup a heartbeat handler
// as a fast-pass. This is to avoid head-of-line blocking from
// disk IO. If Transport does not support this, it can simply
// ignore the call, and push the heartbeat onto the Consumer channel.
func (t *Transport) SetHeartbeatHandler(handler func(rpc raft.RPC)) {
	t.RWMutex.RLock()
	defer t.RWMutex.RUnlock()
	t.heartbeatFn = handler
}

func (t *Transport) getPeer(target raft.ServerAddress) (raftPB.RaftServiceClient, error) {
	t.RWMutex.Lock()         // Locking here
	defer t.RWMutex.Unlock() // Unlocking after the map access is done

	c, ok := t.clients[target]

	if !ok {
		c = &ClientConn{}
		c.mtx.Lock()
		defer c.mtx.Unlock() // We know that Lock was obtained and can use defer here

		t.clients[target] = c

		if c.conn == nil {
			conn, err := grpc.Dial(string(target), t.dialOptions...)
			if err != nil {
				return nil, err
			}
			c.conn = conn
			c.client = raftPB.NewRaftServiceClient(conn)
		}
	}

	return c.client, nil
}
func isHeartbeat(command interface{}) bool {
	req, ok := command.(*raft.AppendEntriesRequest)
	if !ok {
		return false
	}
	if req == nil {
		return false
	}
	return req.Term != 0 &&
		len(req.Leader) != 0 &&
		req.PrevLogEntry == 0 &&
		req.PrevLogTerm == 0 &&
		len(req.Entries) == 0 &&
		req.LeaderCommitIndex == 0
}

func (af *appendFuture) Error() error {
	<-af.done
	return af.err
}
func (af *appendFuture) Start() time.Time {
	return af.start
}

func (af *appendFuture) Request() *raft.AppendEntriesRequest {
	return af.request
}
func (af *appendFuture) Response() *raft.AppendEntriesResponse {
	return &af.response
}

// AppendEntries is used to add another request to the pipeline.
// The send may block which is an effective form of back-pressure.
func (r *raftPipeline) AppendEntries(req *raft.AppendEntriesRequest, resp *raft.AppendEntriesResponse) (raft.AppendFuture, error) {
	af := &appendFuture{
		start:   time.Now(),
		request: req,
		done:    make(chan struct{}),
	}
	if err := r.stream.Send(encoding.EncodeAppendEntriesRequest(req)); err != nil {
		return nil, err
	}
	select {
	case <-r.stream.Context().Done():
		return nil, r.stream.Context().Err()
	case r.inflightCh <- af:
	default:
		return nil, fmt.Errorf("failed to send request to inflightCh")
	}

	return af, nil
}

// Consumer returns a channel that can be used to consume
// response futures when they are ready.
func (r *raftPipeline) Consumer() <-chan raft.AppendFuture {
	return r.doneCh
}

// Close closes the pipeline and cancels all inflight RPCs
func (r *raftPipeline) Close() error {
	r.cancel()
	r.inflightChMtx.Lock()
	defer r.inflightChMtx.Unlock()
	close(r.inflightCh)
	return nil
}

func (r *raftPipeline) receiver() {
	for af := range r.inflightCh {
		af.processMessage(r)
	}
}

// processMessage processes the appendFuture message.
func (af *appendFuture) processMessage(r *raftPipeline) {
	msg, err := r.stream.Recv()
	if err != nil {
		af.err = err
	} else if msg != nil {
		af.response = *encoding.DecodeAppendEntriesResponse(msg)
	}
	close(af.done)
	r.doneCh <- af
}
