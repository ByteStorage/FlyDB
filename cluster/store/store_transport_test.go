package store

import (
	"bytes"
	"github.com/ByteStorage/FlyDB/config"
	raftPB "github.com/ByteStorage/FlyDB/lib/proto/raft"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"reflect"
	"testing"
	"time"
)

type MockClientConn struct {
	AppendEntriesPipelineFn func() (raftPB.RaftService_AppendEntriesPipelineClient, error)
}

func testGetNewRaft(addr string) (*raft.Raft, error) {
	// setup Raft configuration
	s := store{
		id:   addr,
		opts: config.DefaultOptions,
	}
	conf := s.newDefaultConfig()
	conf.LocalID = raft.ServerID(addr)
	s.conf.LogDataStorage = "inMemory"
	s.conf.SnapshotStorage = "memory"
	// setup Raft communication
	s.conf.LocalAddress = raft.ServerAddress(addr)
	l1 := bufconn.Listen(1024)
	trans1, err := newTransport(s.conf, l1, newDialOption(s.conf))
	if err != nil {
		return nil, err
	}

	// create the snapshot store. This allows the Raft to truncate the log.
	snapshots, err := newSnapshotStore(s.conf)
	if err != nil {
		return nil, err
	}

	// create the log store and stable store
	logStore, err := newRaftLog(s.conf)
	if err != nil {
		return nil, err
	}
	stableStore, err := newStableLog(s.conf)
	if err != nil {
		return nil, err
	}

	// create a new finite state machine
	f := newFSM()

	// instantiate the Raft system

	return raft.NewRaft(conf, f, logStore, stableStore, snapshots, trans1)
}
func testGetNewTransport(addr string) (*Transport, error) {
	// setup Raft configuration
	s := store{
		id:   addr,
		opts: config.DefaultOptions,
	}
	s.conf.LogDataStorage = "inMemory"
	s.conf.SnapshotStorage = "memory"
	// setup Raft communication
	s.conf.LocalAddress = raft.ServerAddress(addr)
	s.conf.HeartbeatTimeout = 200 * time.Millisecond
	l1, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	do := newDialOption(s.conf)
	//do = append(do, grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
	//	return l1.Accept()
	//}))

	return newTransport(s.conf, l1, do)
}
func storeTransportTestAppendEntriesResponseFixture() *raft.AppendEntriesResponse {
	return &raft.AppendEntriesResponse{
		RPCHeader:      storeTransportTestRPCHeaderFixture(),
		Term:           0,
		LastLog:        0,
		Success:        false,
		NoRetryBackoff: false,
	}
}
func storeTransportTestAppendEntriesRequestFixture() *raft.AppendEntriesRequest {
	return &raft.AppendEntriesRequest{
		RPCHeader:         storeTransportTestRPCHeaderFixture(),
		Term:              0,
		Leader:            nil,
		PrevLogEntry:      0,
		PrevLogTerm:       0,
		Entries:           nil,
		LeaderCommitIndex: 0,
	}
}
func storeTransportTestRPCHeaderFixture() raft.RPCHeader {
	return raft.RPCHeader{
		ProtocolVersion: 3,
		ID:              []byte(raft.ServerID("")),
		Addr:            []byte(raft.ServerID("127.0.0.1")),
	}
}

func (m *MockClientConn) AppendEntriesPipeline() (raftPB.RaftService_AppendEntriesPipelineClient, error) {
	return m.AppendEntriesPipelineFn()
}

func TestTransport_AppendEntries(t *testing.T) {
	ts1, err := testGetNewTransport("127.0.0.1:8006")
	assert.NoError(t, err)
	ts2, err := testGetNewTransport("127.0.0.1:8007")
	assert.NoError(t, err)
	defer ts1.Close()
	defer ts2.Close()
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-stop:
				return
			case rpc := <-ts2.Consumer():
				if got, want := rpc.Command.(*raft.AppendEntriesRequest).Leader, []byte{3, 2, 1}; !bytes.Equal(got, want) {
					t.Errorf("request.Leader = %v, want %v", got, want)
				}
				if got, want := rpc.Command.(*raft.AppendEntriesRequest).Entries, []*raft.Log{
					{Type: raft.LogNoop, Extensions: []byte{1}, Data: []byte{55}},
				}; !reflect.DeepEqual(got, want) {
					t.Errorf("request.Entries = %v, want %v", got, want)
				}
				rpc.Respond(&raft.AppendEntriesResponse{
					Success: true,
					LastLog: 12396,
				}, nil)
			}
		}
	}()

	var resp raft.AppendEntriesResponse
	if err := ts1.AppendEntries("127.0.0.1:8007", "127.0.0.1:8007", &raft.AppendEntriesRequest{
		RPCHeader: raft.RPCHeader{
			ProtocolVersion: 3,
			ID:              []byte("127.0.0.1:8006"),
			Addr:            []byte("127.0.0.1:8006"),
		},
		Leader: []byte{3, 2, 1},
		Entries: []*raft.Log{
			{Type: raft.LogNoop, Extensions: []byte{1}, Data: []byte{55}},
		},
	}, &resp); err != nil {
		t.Errorf("AppendEntries() failed: %v", err)
	}
	if got, want := resp.LastLog, uint64(12396); got != want {
		t.Errorf("resp.LastLog = %v, want %v", got, want)
	}
	close(stop)
}
func TestTransport_RequestVote(t *testing.T) {
	ts1, err := testGetNewTransport("127.0.0.1:8006")
	assert.NoError(t, err)
	ts2, err := testGetNewTransport("127.0.0.1:8007")
	assert.NoError(t, err)
	defer ts1.Close()
	defer ts2.Close()
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-stop:
				return
			case rpc := <-ts2.Consumer():
				if got, want := rpc.Command.(*raft.RequestVoteRequest).LastLogIndex, uint64(11); !assert.Equal(t, got, want) {
					t.Errorf("request.Leader = %v, want %v", got, want)
				}
				if got, want := rpc.Command.(*raft.RequestVoteRequest).LastLogTerm, uint64(12); !reflect.DeepEqual(got, want) {
					t.Errorf("request.Entries = %v, want %v", got, want)
				}
				rpc.Respond(&raft.RequestVoteResponse{
					RPCHeader: raft.RPCHeader{
						ProtocolVersion: 3,
						ID:              []byte(""),
						Addr:            []byte(""),
					},
					Term:    0,
					Peers:   []byte("hello"),
					Granted: false,
				}, nil)
			}
		}
	}()

	var resp raft.RequestVoteResponse
	if err := ts1.RequestVote("127.0.0.1:8007", "127.0.0.1:8007", &raft.RequestVoteRequest{
		RPCHeader: raft.RPCHeader{
			ProtocolVersion: 3,
			ID:              []byte("127.0.0.1:8006"),
			Addr:            []byte("127.0.0.1:8006"),
		},
		Term:               10,
		Candidate:          []byte("127.0.0.1:8006"),
		LastLogIndex:       11,
		LastLogTerm:        12,
		LeadershipTransfer: false,
	}, &resp); err != nil {
		t.Errorf("RequestVote() failed: %v", err)
	}
	if got, want := resp.Peers, []byte("hello"); !bytes.Equal(got, want) {
		t.Errorf("resp.LastLog = %v, want %v", got, want)
	}
	close(stop)
}
func TestTransport_TimeoutNow(t *testing.T) {
	ts1, err := testGetNewTransport("127.0.0.1:8006")
	assert.NoError(t, err)
	ts2, err := testGetNewTransport("127.0.0.1:8007")
	assert.NoError(t, err)
	defer ts1.Close()
	defer ts2.Close()
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-stop:
				return
			case rpc := <-ts2.Consumer():
				if got, want := rpc.Command.(*raft.TimeoutNowRequest).ID, []byte("127.0.0.1:8006"); !assert.Equal(t, got, want) {
					t.Errorf("request.Leader = %v, want %v", got, want)
				}
				if got, want := rpc.Command.(*raft.TimeoutNowRequest).ProtocolVersion, raft.ProtocolVersion(3); !assert.Equal(t, got, want) {
					t.Errorf("request.Entries = %v, want %v", got, want)
				}
				rpc.Respond(&raft.TimeoutNowResponse{
					RPCHeader: raft.RPCHeader{
						ProtocolVersion: 3,
						ID:              []byte("127.0.0.1:8006"),
						Addr:            []byte("127.0.0.1:8006"),
					},
				}, nil)
			}
		}
	}()

	var resp raft.TimeoutNowResponse
	if err := ts1.TimeoutNow("127.0.0.1:8007", "127.0.0.1:8007", &raft.TimeoutNowRequest{
		RPCHeader: raft.RPCHeader{
			ProtocolVersion: 3,
			ID:              []byte("127.0.0.1:8006"),
			Addr:            []byte("127.0.0.1:8006"),
		},
	}, &resp); err != nil {
		t.Errorf("RequestVote() failed: %v", err)
	}
	if got, want := resp.Addr, []byte("127.0.0.1:8006"); !bytes.Equal(got, want) {
		t.Errorf("resp.LastLog = %v, want %v", got, want)
	}
	close(stop)
}
func TestTransport_Consumer(t *testing.T) {
	// Create a new transport with local address
	conf := config.Config{LocalAddress: "localhost"}
	l1 := bufconn.Listen(1024)
	trans, err := newTransport(conf, l1, newDialOption(conf))
	assert.NoError(t, err)

	c := trans.Consumer()
	assert.NotNil(t, c)
}
func TestTransport_LocalAddr(t *testing.T) {
	// Create a new transport with local address
	conf := config.Config{LocalAddress: "localhost"}
	l1 := bufconn.Listen(1024)
	trans, err := newTransport(conf, l1, newDialOption(conf))
	assert.NoError(t, err)

	c := trans.LocalAddr()
	assert.NotNil(t, c)
}
func TestTransport_EncodePeer(t *testing.T) {
	// Create a new transport with local address
	conf := config.Config{LocalAddress: "localhost"}
	l1 := bufconn.Listen(1024)
	trans, err := newTransport(conf, l1, newDialOption(conf))
	assert.NoError(t, err)

	c := trans.EncodePeer(raft.ServerID("1"), "127.0.0.1")
	assert.NotNil(t, c)
	assert.Equal(t, c, []byte("127.0.0.1"))
}
func TestTransport_DecodePeer(t *testing.T) {
	// Create a new transport with local address
	conf := config.Config{LocalAddress: "localhost"}
	l1 := bufconn.Listen(1024)
	trans, err := newTransport(conf, l1, newDialOption(conf))
	assert.NoError(t, err)

	c := trans.DecodePeer([]byte("127.0.0.1"))
	assert.NotNil(t, c)
	assert.EqualValues(t, c, "127.0.0.1")
}
