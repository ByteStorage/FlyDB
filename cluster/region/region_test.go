package region

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const dirpath = "/tmp/flydb/region_test"

type TestRegionStruct struct {
	region
}

func ReturnNewDB() *engine.DB {
	opts := config.DefaultOptions
	opts.DirPath = dirpath
	db, _ := engine.NewDB(opts)
	return db
}

func destroyRegion(r Region) {
	// Close the region's database
	db := r.(*TestRegionStruct).db
	err := db.Close()
	if err != nil {
		return
	}
	err = os.RemoveAll(dirpath)
	if err != nil {
		return
	}
	r = nil
}

// NewTestRegion creates a new instance of TestRegion.
func NewTestRegion() Region {
	db := ReturnNewDB()
	return &TestRegionStruct{
		region{
			id:         1,
			startKey:   []byte("start"),
			endKey:     []byte("end"),
			db:         db,
			raft:       nil,
			raftGroups: map[uint64]*raft.Raft{},
			leader:     "leader",
			peers:      []string{"peer1", "peer2"},
			size:       100,
		},
	}
}

func TestRegion_Put(t *testing.T) {
	// Create a test region instance
	region := NewTestRegion()
	defer destroyRegion(region)
	err := region.Put([]byte("key_region"), []byte("value_region"))
	assert.NoError(t, err)
}

func TestRegion_Get(t *testing.T) {
	// Create a test region instance
	region := NewTestRegion()
	defer destroyRegion(region)
	err := region.Put([]byte("key_region"), []byte("value_region"))
	assert.NoError(t, err)
	value, err := region.Get([]byte("key_region"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value_region"), value)
}

func TestRegion_Delete(t *testing.T) {
	// Create a test region instance
	region := NewTestRegion()
	defer destroyRegion(region)
	err := region.Put([]byte("key_value"), []byte("value_region"))
	assert.NoError(t, err)
	err = region.Delete([]byte("key_value"))
	assert.NoError(t, err)
	_, err = region.Get([]byte("key_value"))
	assert.Error(t, err)
}

func TestRegion_GetStartKey(t *testing.T) {
	// Create a test region instance
	region := NewTestRegion()
	defer destroyRegion(region)
	startKey := region.GetStartKey()
	assert.Equal(t, []byte("start"), startKey)
}

func TestRegion_GetEndKey(t *testing.T) {
	// Create a test region instance
	region := NewTestRegion()
	defer destroyRegion(region)
	endKey := region.GetEndKey()
	assert.Equal(t, []byte("end"), endKey)
}

func TestRegion_GetLeader(t *testing.T) {
	// Create a test region instance
	region := NewTestRegion()
	defer destroyRegion(region)
	leader := region.GetLeader()
	assert.Equal(t, "leader", leader)
}

func TestRegion_GetPeers(t *testing.T) {
	// Create a test region instance
	region := NewTestRegion()
	defer destroyRegion(region)
	peers := region.GetPeers()
	assert.Equal(t, []string{"peer1", "peer2"}, peers)
}

func TestRegion_GetSize(t *testing.T) {
	// Create a test region instance
	region := NewTestRegion()
	defer destroyRegion(region)
	size := region.GetSize()
	assert.Equal(t, int64(100), size)
}

func TestRegion_AddPeer(t *testing.T) {
	// Create a test region instance
	region := NewTestRegion()
	_ = []struct {
		name          string
		peers         []string
		peerToAdd     string
		expectedPeers []string
		expectError   bool
	}{
		{
			name:          "add a new peer",
			peers:         region.GetPeers(),
			peerToAdd:     "peer3",
			expectedPeers: []string{"peer1", "peer2", "peer3"},
			expectError:   false,
		},
		{
			name:          "add duplicate peer",
			peers:         region.GetPeers(),
			peerToAdd:     "peer1",
			expectedPeers: []string{"peer1", "peer2"},
			expectError:   true,
		},
	}
	destroyRegion(region)

}
