package store

import (
	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
	"io"
)

// transport implements raft.Transport interface
// we can use it to send rpc to other raft nodes
// and receive rpc from other raft nodes
type transport struct {
	//implement me
	localAddr raft.ServerAddress
	consumer  chan raft.RPC
	clients   map[raft.ServerAddress]*grpc.ClientConn
	server    *grpc.Server
}

// NewTransport returns a new transport, it needs start a grpc server
func newTransport() raft.Transport {
	return &transport{}
}

func (t *transport) AppendEntriesPipeline(id raft.ServerID, target raft.ServerAddress) (raft.AppendPipeline, error) {
	//TODO implement me
	panic("implement me")
}

func (t *transport) AppendEntries(id raft.ServerID, target raft.ServerAddress, args *raft.AppendEntriesRequest, resp *raft.AppendEntriesResponse) error {
	//TODO implement me
	panic("implement me")
}

func (t *transport) RequestVote(id raft.ServerID, target raft.ServerAddress, args *raft.RequestVoteRequest, resp *raft.RequestVoteResponse) error {
	//TODO implement me
	panic("implement me")
}

func (t *transport) InstallSnapshot(id raft.ServerID, target raft.ServerAddress, args *raft.InstallSnapshotRequest, resp *raft.InstallSnapshotResponse, data io.Reader) error {
	//TODO implement me
	panic("implement me")
}

func (t *transport) TimeoutNow(id raft.ServerID, target raft.ServerAddress, args *raft.TimeoutNowRequest, resp *raft.TimeoutNowResponse) error {
	//TODO implement me
	panic("implement me")
}

func (t *transport) Consumer() <-chan raft.RPC {
	//implement me
	panic("implement me")
}

func (t *transport) LocalAddr() raft.ServerAddress {
	//implement me
	panic("implement me")
}

func (t *transport) EncodePeer(id raft.ServerID, addr raft.ServerAddress) []byte {
	//implement me
	panic("implement me")
}

func (t *transport) DecodePeer([]byte) raft.ServerAddress {
	//implement me
	panic("implement me")
}

func (t *transport) SetHeartbeatHandler(handler func(rpc raft.RPC)) {
	//implement me
	panic("implement me")
}
