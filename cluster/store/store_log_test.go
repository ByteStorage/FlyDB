package store

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/datastore"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestRegister(t *testing.T) {
	foo := func(conf config.Config) (raft.LogStore, error) {
		return &datastore.InMemStore{}, nil
	}
	opts := config.DefaultOptions
	conf := config.Config{}
	conf.LogDataStorageSize = 100 * 100 * 32
	conf.LogDataStoragePath = opts.DirPath
	err := Register("memory", foo)
	assert.NoError(t, err)
	// get the datastore and check
	conf.LogDataStorage = "memory"
	ds, err := getDataStore(conf)
	assert.NoError(t, err)
	assert.IsType(t, &datastore.InMemStore{}, ds)
}

func TestInit(t *testing.T) {
	// initialize the dbs
	_ = Init()
	// in memory DB
	opts := config.DefaultOptions
	conf := config.Config{}
	conf.LogDataStorageSize = 100 * 100 * 32
	conf.LogDataStoragePath = opts.DirPath
	conf.LogDataStorage = "memory"
	ds, err := getDataStore(conf)
	assert.NoError(t, err)
	assert.IsType(t, &datastore.InMemStore{}, ds)
	// FlyDB
	tf, err := testTempFile()
	assert.NoError(t, err)
	conf.LogDataStoragePath = tf
	conf.LogDataStorage = "flydb"
	ds, err = getDataStore(conf)
	assert.NoError(t, err)
	assert.IsType(t, &datastore.FlyDbStore{}, ds)
	// RockDB
	tf, err = testTempDir()
	assert.NoError(t, err)
	conf.LogDataStoragePath = tf
	conf.LogDataStorage = "rocksdb"
	ds, err = getDataStore(conf)
	assert.NoError(t, err)
	assert.IsType(t, &datastore.RocksDbStore{}, ds)
	// BoltDB
	tf, err = testTempFile()
	assert.NoError(t, err)
	conf.LogDataStorage = "boltdb"
	conf.LogDataStoragePath = tf
	ds, err = getDataStore(conf)
	assert.NoError(t, err)
	assert.IsType(t, &datastore.BoltDbStore{}, ds)

}

func testTempFile() (string, error) {
	tmpFile, err := os.CreateTemp("", "test_flydb_storage")
	if err != nil {
		return "", err
	}
	err = os.Remove(tmpFile.Name())
	if err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}
func testTempDir() (string, error) {
	tmpDir, err := os.MkdirTemp("", "test_flydb_rocksDB")
	if err != nil {
		return "", err
	}
	err = os.Remove(tmpDir)
	if err != nil {
		return "", err
	}
	return tmpDir, nil
}
