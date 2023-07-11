package store

import (
	"github.com/hashicorp/raft"
	"os"
	"time"
)

// newDefaultConfig returns a new default raft config
func (s *store) newDefaultConfig() *raft.Config {
	return &raft.Config{
		ProtocolVersion:          raft.ProtocolVersionMax, // using latest protocol version
		HeartbeatTimeout:         1000 * time.Millisecond,
		ElectionTimeout:          1000 * time.Millisecond,
		CommitTimeout:            50 * time.Millisecond,
		MaxAppendEntries:         64,
		BatchApplyCh:             false,
		ShutdownOnRemove:         true,
		TrailingLogs:             10240,
		SnapshotInterval:         120 * time.Second,
		SnapshotThreshold:        8192,
		LeaderLeaseTimeout:       500 * time.Millisecond,
		LogOutput:                os.Stderr,
		LogLevel:                 "debug", // this is for `dev` environments, for production, use `info` or `warn`.
		Logger:                   nil,     // this by default uses hclog which we don't use
		NoSnapshotRestoreOnStart: false,
	}
}
