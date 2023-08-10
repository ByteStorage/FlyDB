package meta

// create a simple raft node for the meta node
func (m *meta) newRaftNode() {
	// All new methods below can add other return values as needed, such as err

	// create default config for raft
	//raftConfig := raft.DefaultConfig()

	// setup Raft communication

	// create the snapshot store. This allows the Raft to truncate the log.

	// create the log store and stable store

	// create a new finite state machine

	// instantiate the Raft system
	//r, err := raft.NewRaft(raftConfig, f, logStore, stableStore, snapshots, t)
	//if err != nil {
	//	return nil, err
	//}
	panic("implement me")
}
