package store

import (
	"github.com/ByteStorage/FlyDB/config"
	raftPB "github.com/ByteStorage/FlyDB/lib/proto/raft"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockClientConn struct {
	AppendEntriesPipelineFn func() (raftPB.RaftService_AppendEntriesPipelineClient, error)
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
func TestNewTransport(t *testing.T) {
	conf := config.Config{
		LocalAddress:     "localhost",
		HeartbeatTimeout: 1000,
	}

	transport := newTransport(conf)

	assert.Equal(t, conf.LocalAddress, transport.LocalAddr())
	assert.NotNil(t, transport.Consumer())
}

func TestTransport_AppendEntriesPipeline(t *testing.T) {
	// Create a new transport with local address
	//conf := config.Config{LocalAddress: "localhost"}
	//trans := newTransport(conf)

	//_, err := trans.AppendEntriesPipeline("id", "localhost:100")
	//assert.Nil(t, err)
	// send some data should not return error
}
func TestTransport_Consumer(t *testing.T) {
	// Create a new transport with local address
	conf := config.Config{LocalAddress: "localhost"}
	trans := newTransport(conf)

	c := trans.Consumer()
	assert.NotNil(t, c)
}
func TestTransport_LocalAddr(t *testing.T) {
	// Create a new transport with local address
	conf := config.Config{LocalAddress: "localhost"}
	trans := newTransport(conf)

	c := trans.LocalAddr()
	assert.NotNil(t, c)
}
func TestTransport_EncodePeer(t *testing.T) {
	// Create a new transport with local address
	conf := config.Config{LocalAddress: "localhost"}
	trans := newTransport(conf)

	c := trans.EncodePeer(raft.ServerID("1"), "127.0.0.1")
	assert.NotNil(t, c)
	assert.Equal(t, c, []byte("127.0.0.1"))
}
func TestTransport_DecodePeer(t *testing.T) {
	// Create a new transport with local address
	conf := config.Config{LocalAddress: "localhost"}
	trans := newTransport(conf)

	c := trans.DecodePeer([]byte("127.0.0.1"))
	assert.NotNil(t, c)
	assert.EqualValues(t, c, "127.0.0.1")
}
