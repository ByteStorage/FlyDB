package store

import (
	"github.com/hashicorp/raft"
)

// newDefaultConfig returns a new default raft config
func (s *store) newDefaultConfig() *raft.Config {
	return &raft.Config{
		// implement me
	}
}

// you should read this to ensure suitable config for FlyDB raft, and fill in the newDefaultConfig
// You can check the method `raft.DefaultConfig()` of raft

// // Config provides any necessary configuration for the Raft server.
//type Config struct {
//	// ProtocolVersion allows a Raft server to inter-operate with older
//	// Raft servers running an older version of the code. This is used to
//	// version the wire protocol as well as Raft-specific log entries that
//	// the server uses when _speaking_ to other servers. There is currently
//	// no auto-negotiation of versions so all servers must be manually
//	// configured with compatible versions. See ProtocolVersionMin and
//	// ProtocolVersionMax for the versions of the protocol that this server
//	// can _understand_.
//	ProtocolVersion ProtocolVersion
//
//	// HeartbeatTimeout specifies the time in follower state without contact
//	// from a leader before we attempt an election.
//	HeartbeatTimeout time.Duration
//
//	// ElectionTimeout specifies the time in candidate state without contact
//	// from a leader before we attempt an election.
//	ElectionTimeout time.Duration
//
//	// CommitTimeout specifies the time without an Apply operation before the
//	// leader sends an AppendEntry RPC to followers, to ensure a timely commit of
//	// log entries.
//	// Due to random staggering, may be delayed as much as 2x this value.
//	CommitTimeout time.Duration
//
//	// MaxAppendEntries controls the maximum number of append entries
//	// to send at once. We want to strike a balance between efficiency
//	// and avoiding waste if the follower is going to reject because of
//	// an inconsistent log.
//	MaxAppendEntries int
//
//	// BatchApplyCh indicates whether we should buffer applyCh
//	// to size MaxAppendEntries. This enables batch log commitment,
//	// but breaks the timeout guarantee on Apply. Specifically,
//	// a log can be added to the applyCh buffer but not actually be
//	// processed until after the specified timeout.
//	BatchApplyCh bool
//
//	// If we are a member of a cluster, and RemovePeer is invoked for the
//	// local node, then we forget all peers and transition into the follower state.
//	// If ShutdownOnRemove is set, we additional shutdown Raft. Otherwise,
//	// we can become a leader of a cluster containing only this node.
//	ShutdownOnRemove bool
//
//	// TrailingLogs controls how many logs we leave after a snapshot. This is used
//	// so that we can quickly replay logs on a follower instead of being forced to
//	// send an entire snapshot. The value passed here is the initial setting used.
//	// This can be tuned during operation using ReloadConfig.
//	TrailingLogs uint64
//
//	// SnapshotInterval controls how often we check if we should perform a
//	// snapshot. We randomly stagger between this value and 2x this value to avoid
//	// the entire cluster from performing a snapshot at once. The value passed
//	// here is the initial setting used. This can be tuned during operation using
//	// ReloadConfig.
//	SnapshotInterval time.Duration
//
//	// SnapshotThreshold controls how many outstanding logs there must be before
//	// we perform a snapshot. This is to prevent excessive snapshotting by
//	// replaying a small set of logs instead. The value passed here is the initial
//	// setting used. This can be tuned during operation using ReloadConfig.
//	SnapshotThreshold uint64
//
//	// LeaderLeaseTimeout is used to control how long the "lease" lasts
//	// for being the leader without being able to contact a quorum
//	// of nodes. If we reach this interval without contact, we will
//	// step down as leader.
//	LeaderLeaseTimeout time.Duration
//
//	// LocalID is a unique ID for this server across all time. When running with
//	// ProtocolVersion < 3, you must set this to be the same as the network
//	// address of your transport.
//	LocalID ServerID
//
//	// NotifyCh is used to provide a channel that will be notified of leadership
//	// changes. Raft will block writing to this channel, so it should either be
//	// buffered or aggressively consumed.
//	NotifyCh chan<- bool
//
//	// LogOutput is used as a sink for logs, unless Logger is specified.
//	// Defaults to os.Stderr.
//	LogOutput io.Writer
//
//	// LogLevel represents a log level. If the value does not match a known
//	// logging level hclog.NoLevel is used.
//	LogLevel string
//
//	// Logger is a user-provided logger. If nil, a logger writing to
//	// LogOutput with LogLevel is used.
//	Logger hclog.Logger
//
//	// NoSnapshotRestoreOnStart controls if raft will restore a snapshot to the
//	// FSM on start. This is useful if your FSM recovers from other mechanisms
//	// than raft snapshotting. Snapshot metadata will still be used to initialize
//	// raft's configuration and index values.
//	NoSnapshotRestoreOnStart bool
//
//	// skipStartup allows NewRaft() to bypass all background work goroutines
//	skipStartup bool
//}
