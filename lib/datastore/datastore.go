package datastore

import "github.com/hashicorp/raft"

type DataStore interface {
	raft.LogStore
	raft.StableStore
}
