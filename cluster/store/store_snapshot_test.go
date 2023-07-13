package store

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestRegisterSnapshotStore(t *testing.T) {
	type test struct {
		name        string
		datastore   string
		want        raft.SnapshotStore
		expectError bool
	}
	tests := []test{
		{
			name:        "memory",
			datastore:   "memory",
			want:        &raft.InmemSnapshotStore{},
			expectError: false,
		},
		{
			name:        "correctly get file",
			datastore:   "file",
			want:        &raft.FileSnapshotStore{},
			expectError: false,
		},
		{
			name:        "correctly get discard",
			datastore:   "discard",
			want:        &raft.DiscardSnapshotStore{},
			expectError: false,
		},
		{
			name:        "a wrong datastore",
			datastore:   "wrong",
			want:        nil,
			expectError: true,
		},
	}
	for _, tc := range tests {
		conf := config.Config{}
		conf.SnapshotStorage = tc.datastore
		conf.SnapshotStoragePath, _ = testTempFile()
		snst, err := getSnapShotStore(conf)
		if tc.expectError {
			assert.NotNil(t, err)
		} else {
			assert.IsType(t, tc.want, snst)
			checkStoreSnapshotInterfaceEmbedding(t, snst)
		}
	}
}

func checkStoreSnapshotInterfaceEmbedding(t *testing.T, in interface{}) {
	assert.True(t, reflect.TypeOf(in).Implements(reflect.TypeOf((*raft.SnapshotStore)(nil)).Elem()))
}
