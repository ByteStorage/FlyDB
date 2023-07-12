package datastore

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func testRocksDbDatastore() (raft.LogStore, error) {
	//fmt.Sprintf("%s/%s", os.TempDir, "")
	tmpDir, err := os.MkdirTemp("", "test_flydb_rocksDB")
	if err != nil {
		return nil, err
	}
	err = os.Remove(tmpDir)
	if err != nil {
		return nil, err
	}

	// Successfully creates and returns a store

	return NewLogRocksDbStorage(config.Config{LogDataStoragePath: tmpDir})
}

func TestRocksDbStore_DeleteRange(t *testing.T) {
	r, err := testRocksDbDatastore()
	assert.NoError(t, err)
	logs := []*raft.Log{
		createRaftLog(1, "test1"),
		createRaftLog(2, "test2"),
		createRaftLog(23, "test3"),
		createRaftLog(6, "test4"),
	}
	err = r.StoreLogs(logs)
	assert.NoError(t, err)
	var log1, log2, log3, log4 raft.Log
	_ = r.GetLog(23, &log1)
	assert.Equal(t, *logs[2], log1)

	_ = r.DeleteRange(9, 50)
	assert.NoError(t, err)
	_ = r.GetLog(23, &log2)
	assert.Equal(t, raft.Log{}, log2)
	// delete another range
	_ = r.GetLog(2, &log3)
	assert.Equal(t, *logs[1], log3)

	_ = r.DeleteRange(2, 5)
	assert.NoError(t, err)
	_ = r.GetLog(2, &log4)
	assert.Equal(t, raft.Log{}, log4)
}
func TestRocksDbStore_FirstIndex(t *testing.T) {
	r, err := testRocksDbDatastore()
	assert.NoError(t, err)
	logs := []*raft.Log{
		createRaftLog(8, "test2"),
		createRaftLog(1, "test1"),
		createRaftLog(23, "test3"),
	}
	err = r.StoreLogs(logs)
	assert.NoError(t, err)

	fi, err := r.FirstIndex()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, fi)
	// remove first element
	_ = r.DeleteRange(1, 2)
	fi, err = r.FirstIndex()
	assert.NoError(t, err)
	assert.EqualValues(t, 8, fi)

}
func TestRocksDbStore_LastIndex(t *testing.T) {
	r, err := testRocksDbDatastore()
	assert.NoError(t, err)
	logs := []*raft.Log{
		createRaftLog(1, "test1"),
		createRaftLog(2, "test2"),
		createRaftLog(23, "test3"),
		createRaftLog(6, "test4"),
	}
	err = r.StoreLogs(logs)
	assert.NoError(t, err)

	//
	li, err := r.LastIndex()
	assert.NoError(t, err)
	assert.EqualValues(t, 23, li)
	// remove first element
	_ = r.DeleteRange(9, 50)
	li, err = r.LastIndex()
	assert.NoError(t, err)
	assert.EqualValues(t, 6, li)

}
func TestRocksDbStore_GetLog(t *testing.T) {
	r, err := testRocksDbDatastore()
	assert.NoError(t, err)
	logs := []*raft.Log{
		createRaftLog(1, "test1"),
		createRaftLog(2, "test2"),
		createRaftLog(23, "test3"),
	}
	err = r.StoreLogs(logs)
	assert.NoError(t, err)
	// get the logs
	var log raft.Log
	err = r.GetLog(1, &log)
	assert.NoError(t, err)
	assert.Equal(t, log, *logs[0])
	// get log of id 23
	err = r.GetLog(23, &log)
	assert.NoError(t, err)
	assert.Equal(t, log, *logs[2])

}
func TestRocksDbStore_StoreLog(t *testing.T) {
	r, err := testRocksDbDatastore()
	assert.NoError(t, err)

	l1 := createRaftLog(23, "test4")
	err = r.StoreLog(l1)
	assert.NoError(t, err)
	// get the logs
	var log1, log2, log3, log5 raft.Log
	_ = r.GetLog(1, &log1)
	_ = r.GetLog(23, &log3)
	assert.NoError(t, err)

	assert.Equal(t, log1, raft.Log{})
	assert.Equal(t, log2, raft.Log{})
	assert.Equal(t, log3, *l1)

	// check for errors
	err = r.GetLog(4, &log5)
	assert.True(t, err != nil)

}
func TestRocksDbStore_StoreLogs(t *testing.T) {
	r, err := testRocksDbDatastore()
	assert.NoError(t, err)

	logs := []*raft.Log{
		createRaftLog(1, "test1"),
		createRaftLog(23, "test3"),
		createRaftLog(8, "test2"),
	}
	err = r.StoreLogs(logs)
	assert.NoError(t, err)
	// get the logs
	var log1, log2, log3, log5 raft.Log
	_ = r.GetLog(1, &log1)
	_ = r.GetLog(23, &log2)
	_ = r.GetLog(8, &log3)

	assert.Equal(t, log1, *logs[0])
	assert.Equal(t, log2, *logs[1])
	assert.Equal(t, log3, *logs[2])

	// check for errors
	err = r.GetLog(4, &log5)
	assert.True(t, err != nil)

}
