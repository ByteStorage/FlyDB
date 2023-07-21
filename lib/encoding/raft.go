package encoding

import (
	raftPB "github.com/ByteStorage/FlyDB/lib/proto/raft"
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EncodeRPCHeader(s raft.RPCHeader) *raftPB.RPCHeader {
	return &raftPB.RPCHeader{
		ProtocolVersion: int64(s.ProtocolVersion),
		ID:              s.ID,
		Addr:            s.Addr,
	}
}
func DecodeRPCHeader(m *raftPB.RPCHeader) raft.RPCHeader {
	return raft.RPCHeader{
		ProtocolVersion: raft.ProtocolVersion(m.ProtocolVersion),
		ID:              m.ID,
		Addr:            m.Addr,
	}
}

func EncodeRequestVoteRequest(s *raft.RequestVoteRequest) *raftPB.RequestVoteRequest {
	return &raftPB.RequestVoteRequest{
		RpcHeader:          EncodeRPCHeader(s.RPCHeader),
		Term:               s.Term,
		Candidate:          s.Candidate,
		LastLogIndex:       s.LastLogIndex,
		LastLogTerm:        s.LastLogTerm,
		LeadershipTransfer: s.LeadershipTransfer,
	}
}

func DecodeRequestVoteResponse(m *raftPB.RequestVoteResponse) *raft.RequestVoteResponse {
	return &raft.RequestVoteResponse{
		RPCHeader: DecodeRPCHeader(m.RpcHeader),
		Term:      m.Term,
		Peers:     m.Peers,
		Granted:   m.Granted,
	}
}

func EncodeAppendEntriesRequest(s *raft.AppendEntriesRequest) *raftPB.AppendEntriesRequest {
	return &raftPB.AppendEntriesRequest{
		RpcHeader:         EncodeRPCHeader(s.RPCHeader),
		Term:              s.Term,
		Leader:            s.Leader,
		PrevLogEntry:      s.PrevLogEntry,
		PrevLogTerm:       s.PrevLogTerm,
		Entries:           encodeLogs(s.Entries),
		LeaderCommitIndex: s.LeaderCommitIndex,
	}
}
func DecodeAppendEntriesResponse(m *raftPB.AppendEntriesResponse) *raft.AppendEntriesResponse {
	return &raft.AppendEntriesResponse{
		RPCHeader:      DecodeRPCHeader(m.RpcHeader),
		Term:           m.Term,
		LastLog:        m.LastLog,
		Success:        m.Success,
		NoRetryBackoff: m.NoRetryBackoff,
	}
}
func EncodeTimeoutNowRequest(s *raft.TimeoutNowRequest) *raftPB.TimeoutNowRequest {
	return &raftPB.TimeoutNowRequest{
		RpcHeader: EncodeRPCHeader(s.RPCHeader),
	}
}
func DecodeTimeoutNowResponse(m *raftPB.TimeoutNowResponse) *raft.TimeoutNowResponse {
	return &raft.TimeoutNowResponse{
		RPCHeader: DecodeRPCHeader(m.RpcHeader),
	}
}
func encodeLogs(s []*raft.Log) []*raftPB.Log {
	ret := make([]*raftPB.Log, len(s))
	for i, l := range s {
		ret[i] = encodeLog(l)
	}
	return ret
}

func encodeLog(l *raft.Log) *raftPB.Log {
	ret := &raftPB.Log{
		Index:      l.Index,
		Term:       l.Term,
		Type:       encodeLogType(l.Type),
		Data:       l.Data,
		Extensions: l.Extensions,
		AppendedAt: timestamppb.New(l.AppendedAt),
	}
	return ret
}
func encodeLogType(s raft.LogType) raftPB.Log_LogType {
	switch s {
	case raft.LogCommand:
		return raftPB.Log_LOG_COMMAND
	case raft.LogNoop:
		return raftPB.Log_LOG_NOOP
	case raft.LogAddPeerDeprecated:
		return raftPB.Log_LOG_ADD_PEER_DEPRECATED
	case raft.LogRemovePeerDeprecated:
		return raftPB.Log_LOG_REMOVE_PEER_DEPRECATED
	case raft.LogBarrier:
		return raftPB.Log_LOG_BARRIER
	case raft.LogConfiguration:
		return raftPB.Log_LOG_CONFIGURATION
	default:
		panic("invalid LogType")
	}
}
func EncodeInstallSnapshotRequest(s *raft.InstallSnapshotRequest) *raftPB.InstallSnapshotRequest {
	return &raftPB.InstallSnapshotRequest{
		RpcHeader:          EncodeRPCHeader(s.RPCHeader),
		SnapshotVersion:    int64(s.SnapshotVersion),
		Term:               s.Term,
		Leader:             s.Leader,
		LastLogIndex:       s.LastLogIndex,
		LastLogTerm:        s.LastLogTerm,
		Peers:              s.Peers,
		Configuration:      s.Configuration,
		ConfigurationIndex: s.ConfigurationIndex,
		Size:               s.Size,
	}
}
func DecodeInstallSnapshotResponse(m *raftPB.InstallSnapshotResponse) *raft.InstallSnapshotResponse {
	return &raft.InstallSnapshotResponse{
		RPCHeader: DecodeRPCHeader(m.RpcHeader),
		Term:      m.Term,
		Success:   m.Success,
	}
}
