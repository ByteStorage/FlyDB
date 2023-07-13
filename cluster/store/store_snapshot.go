package store

import (
	"errors"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/hashicorp/raft"
	"os"
	"path"
)

func newSnapshotStore(conf config.Config) (raft.SnapshotStore, error) {
	snStore, err := getSnapShotStore(conf)
	if err != nil {
		return nil, err
	}
	return snStore, nil
}

func getSnapShotStore(conf config.Config) (raft.SnapshotStore, error) {
	snapshotStoreDir := path.Join(conf.SnapshotStoragePath, "snapshot")

	switch conf.SnapshotStorage {
	case "memory":
		return raft.NewInmemSnapshotStore(), nil
	case "discard":
		return raft.NewDiscardSnapshotStore(), nil
	case "file":
		return raft.NewFileSnapshotStore(snapshotStoreDir, 5, os.Stderr) //todo get the retain param from config
	default:
		return nil, errors.New("the datastore does not exist")
	}
}
