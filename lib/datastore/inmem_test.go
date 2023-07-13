package datastore

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testMemoryDatastore() raft.LogStore {
	ds, _ := NewLogInMemStorage(config.Config{})
	return ds
}
func createRaftLog(idx uint64, data string) *raft.Log {
	return &raft.Log{
		Data:  []byte(data),
		Index: idx,
	}
}

func TestInMemStore_GetLog(t *testing.T) {
	ds := testMemoryDatastore()
	logs := []*raft.Log{
		createRaftLog(1, "test1"),
		createRaftLog(2, "test2"),
		createRaftLog(23, "test3"),
	}
	err := ds.StoreLogs(logs)
	assert.NoError(t, err)
	// get the logs
	var log raft.Log
	err = ds.GetLog(1, &log)
	assert.NoError(t, err)
	assert.Equal(t, log, *logs[0])
	// get log of id 23
	err = ds.GetLog(23, &log)
	assert.NoError(t, err)
	assert.Equal(t, log, *logs[2])
}

func TestInMemStore_DeleteRange(t *testing.T) {
	ds := testMemoryDatastore()
	logs := []*raft.Log{
		createRaftLog(1, "test1"),
		createRaftLog(2, "test2"),
		createRaftLog(23, "test3"),
		createRaftLog(6, "test4"),
	}
	err := ds.StoreLogs(logs)
	assert.NoError(t, err)
	var log1, log2 raft.Log
	_ = ds.GetLog(23, &log1)
	assert.Equal(t, *logs[2], log1)

	_ = ds.DeleteRange(9, 50)
	assert.NoError(t, err)
	_ = ds.GetLog(23, &log2)
	assert.Equal(t, raft.Log{}, log2)

}

func TestInMemStore_StoreLog(t *testing.T) {
	ds := testMemoryDatastore()

	l1 := createRaftLog(6, "test4")
	err := ds.StoreLog(l1)
	assert.NoError(t, err)
	// get the logs
	var log1, log2, log3, log4, log5 raft.Log
	_ = ds.GetLog(1, &log1)
	_ = ds.GetLog(2, &log2)
	_ = ds.GetLog(23, &log3)
	err = ds.GetLog(6, &log4)
	assert.NoError(t, err)

	assert.Equal(t, log1, raft.Log{})
	assert.Equal(t, log2, raft.Log{})
	assert.Equal(t, log3, raft.Log{})
	assert.Equal(t, log4, *l1)

	// check for errors
	err = ds.GetLog(4, &log5)
	assert.True(t, err != nil)

}

func TestInMemStore_StoreLogs(t *testing.T) {
	ds := testMemoryDatastore()
	logs := []*raft.Log{
		createRaftLog(1, "test1"),
		createRaftLog(2, "test2"),
		createRaftLog(23, "test3"),
		createRaftLog(6, "test4"),
	}
	err := ds.StoreLogs(logs)
	assert.NoError(t, err)
	// get the logs
	var log1, log2, log3, log4, log5 raft.Log
	_ = ds.GetLog(1, &log1)
	_ = ds.GetLog(2, &log2)
	_ = ds.GetLog(23, &log3)
	_ = ds.GetLog(6, &log4)

	assert.Equal(t, log1, *logs[0])
	assert.Equal(t, log2, *logs[1])
	assert.Equal(t, log3, *logs[2])
	assert.Equal(t, log4, *logs[3])

	// check for errors
	err = ds.GetLog(4, &log5)
	assert.True(t, err != nil)

}

func TestInMemStore_FirstIndex(t *testing.T) {
	ds := testMemoryDatastore()
	logs := []*raft.Log{
		createRaftLog(8, "test2"),
		createRaftLog(1, "test1"),
		createRaftLog(23, "test3"),
	}
	err := ds.StoreLogs(logs)
	assert.NoError(t, err)

	fi, err := ds.FirstIndex()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, fi)
	// remove first element
	_ = ds.DeleteRange(1, 2)
	fi, err = ds.FirstIndex()
	assert.NoError(t, err)
	assert.EqualValues(t, 8, fi)

}

func TestInMemStore_LastIndex(t *testing.T) {
	ds := testMemoryDatastore()
	logs := []*raft.Log{
		createRaftLog(1, "test1"),
		createRaftLog(2, "test2"),
		createRaftLog(23, "test3"),
		createRaftLog(6, "test4"),
	}
	err := ds.StoreLogs(logs)
	assert.NoError(t, err)

	li, err := ds.LastIndex()
	assert.NoError(t, err)
	assert.EqualValues(t, 23, li)
	// remove first element
	_ = ds.DeleteRange(9, 50)
	li, err = ds.LastIndex()
	assert.NoError(t, err)
	assert.EqualValues(t, 6, li)

}
