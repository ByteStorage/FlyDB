package datastore

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func testBoltDatastore() (DataStore, error) {
	tmpFile, err := os.CreateTemp("", "test_flydb_boltdb")
	if err != nil {
		return nil, err
	}
	err = os.Remove(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	// Successfully creates and returns a store

	return NewLogBoltDbStorage(config.Config{LogDataStoragePath: tmpFile.Name()})
}

func TestBoltDbStore_DeleteRange(t *testing.T) {
	r, err := testBoltDatastore()
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
func TestBoltDbStore_FirstIndex(t *testing.T) {
	r, err := testBoltDatastore()
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
func TestBoltDbStore_LastIndex(t *testing.T) {
	r, err := testBoltDatastore()
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
func TestBoltDbStore_GetLog(t *testing.T) {
	r, err := testBoltDatastore()
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
func TestBoltDbStore_StoreLog(t *testing.T) {
	r, err := testBoltDatastore()
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
func TestBoltDbStore_StoreLogs(t *testing.T) {
	r, err := testBoltDatastore()
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

func TestBoltDbStore_Set(t *testing.T) {
	ds, err := testBoltDatastore()
	assert.NoError(t, err)
	type kv struct {
		key string
		val string
	}
	type test struct {
		input       []kv
		expectError bool
	}
	tests := []test{
		{
			input: []kv{
				{key: "1", val: "2"},
				{key: "foo", val: "bar"},
				{key: "hello", val: "world"},
			},
			expectError: false,
		},
		{
			input: []kv{
				{key: "", val: "bar"},
			},
			expectError: true,
		},
	}
	for _, tc := range tests {
		// set all inputs
		for _, v := range tc.input {
			err := ds.Set(stringToBytes(v.key), stringToBytes(v.val))
			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
		}
		// recall all inputs
		for _, v := range tc.input {
			val, err := ds.Get(stringToBytes(v.key))
			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, v.val, string(val))
			}
		}

	}

}
func TestBoltDbStore_Get(t *testing.T) {
	type kv struct {
		key string
		val string
	}
	type test struct {
		name        string
		input       []kv
		query       []kv
		expectError bool
	}
	tests := []test{
		{
			name: "set three",
			input: []kv{
				{key: "1", val: "2"},
				{key: "foo", val: "bar"},
				{key: "hello", val: "world"},
			},
			query: []kv{
				{key: "1", val: "2"},
				{key: "foo", val: "bar"},
				{key: "hello", val: "world"},
			},
			expectError: false,
		},
		{
			name: "non existence",
			input: []kv{
				{key: "4", val: "bar"},
			},
			query: []kv{
				{key: "1", val: ""},
				{key: "2", val: ""},
			},
			expectError: true,
		},
	}
	for _, tc := range tests {
		ds, err := testBoltDatastore()
		assert.NoError(t, err)
		// set all inputs
		for _, v := range tc.input {
			_ = ds.Set(stringToBytes(v.key), stringToBytes(v.val))

		}
		// recall all inputs
		for _, v := range tc.query {
			val, err := ds.Get(stringToBytes(v.key))
			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, v.val, string(val))
			}
		}

	}

}
func TestBoltDbStore_SetUint64(t *testing.T) {
	ds, err := testBoltDatastore()
	assert.NoError(t, err)
	type kv struct {
		key string
		val uint64
	}
	type test struct {
		input       []kv
		expectError bool
	}
	tests := []test{
		{
			input: []kv{
				{key: "1", val: 2343},
				{key: "foo", val: 23},
				{key: "hello", val: 5645},
			},
			expectError: false,
		},
		{
			input: []kv{
				{key: "", val: 654},
			},
			expectError: true,
		},
	}
	for _, tc := range tests {
		// set all inputs
		for _, v := range tc.input {
			err := ds.Set(stringToBytes(v.key), uint64ToBytes(v.val))
			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
		}
		// recall all inputs
		for _, v := range tc.input {
			val, err := ds.Get(stringToBytes(v.key))
			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, v.val, bytesToUint64(val))
			}
		}
	}

}
func TestBoltDbStore_GetUint64(t *testing.T) {
	type kv struct {
		key string
		val uint64
	}
	type test struct {
		name        string
		input       []kv
		query       []kv
		expectError bool
	}
	tests := []test{
		{
			name: "set three",
			input: []kv{
				{key: "1", val: 11},
				{key: "foo", val: 55},
				{key: "hello", val: 336},
			},
			query: []kv{
				{key: "1", val: 11},
				{key: "foo", val: 55},
				{key: "hello", val: 336},
			},
			expectError: false,
		},
		{
			name: "non existence",
			input: []kv{
				{key: "4", val: 2},
			},
			query: []kv{
				{key: "1", val: 2},
				{key: "2", val: 4},
			},
			expectError: true,
		},
	}
	for _, tc := range tests {
		ds, err := testBoltDatastore()
		assert.NoError(t, err)
		// set all inputs
		for _, v := range tc.input {
			_ = ds.Set(stringToBytes(v.key), uint64ToBytes(v.val))

		}
		// recall all inputs
		for _, v := range tc.query {
			val, err := ds.Get(stringToBytes(v.key))
			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, v.val, bytesToUint64(val))
			}
		}

	}

}
