package structure

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/db/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
)

// DataStructureManager provides unified access to all data structures
// Supports all data structure types through a single DB instance
type DataStructureManager struct {
	db             *engine.DB
	stringStruct   *StringStructureInternal
	hashStruct     *HashStructureInternal
	listStruct     *ListStructureInternal
	setStruct      *SetStructureInternal
	zsetStruct     *ZSetStructureInternal
	bitmapStruct   *BitmapStructureInternal
	streamStruct   *StreamStructureInternal
	expiringStruct *ExpiringKeyInternal
}

// NewDataStructureManager creates a new data structure manager
// It initializes all supported data structures sharing a single DB instance
func NewDataStructureManager(options config.Options) (*DataStructureManager, error) {
	// Create a shared DB instance
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}

	// Create manager instance
	manager := &DataStructureManager{
		db: db,
	}

	// Initialize each data structure, sharing the same DB instance
	manager.stringStruct = newStringStructureInternal(db)
	manager.hashStruct = newHashStructureInternal(db)
	manager.listStruct = newListStructureInternal(db)
	manager.setStruct = newSetStructureInternal(db)
	manager.zsetStruct = newZSetStructureInternal(db)
	manager.bitmapStruct = newBitmapStructureInternal(db)
	manager.streamStruct = newStreamStructureInternal(db)
	manager.expiringStruct = newExpiringKeyInternal(db)

	return manager, nil
}

// String returns the string data structure operation interface
func (m *DataStructureManager) String() *StringStructureInternal {
	return m.stringStruct
}

// Hash returns the hash data structure operation interface
func (m *DataStructureManager) Hash() *HashStructureInternal {
	return m.hashStruct
}

// List returns the list data structure operation interface
func (m *DataStructureManager) List() *ListStructureInternal {
	return m.listStruct
}

// Set returns the set data structure operation interface
func (m *DataStructureManager) Set() *SetStructureInternal {
	return m.setStruct
}

// ZSet returns the sorted set data structure operation interface
func (m *DataStructureManager) ZSet() *ZSetStructureInternal {
	return m.zsetStruct
}

// Bitmap returns the bitmap data structure operation interface
func (m *DataStructureManager) Bitmap() *BitmapStructureInternal {
	return m.bitmapStruct
}

// Stream returns the stream data structure operation interface
func (m *DataStructureManager) Stream() *StreamStructureInternal {
	return m.streamStruct
}

// Keys returns all keys matching the given pattern
func (m *DataStructureManager) Keys(pattern string) ([]string, error) {
	// Need to implement a method to find keys from all data structures
	// Temporarily simplified implementation using StringStructure's Keys method
	return m.stringStruct.Keys(pattern)
}

// Type returns the data structure type of the key
func (m *DataStructureManager) Type(key string) (string, error) {
	// Check string type
	if exists, _ := m.stringStruct.Exists(key); exists {
		return "string", nil
	}

	// Check hash type
	// Need to implement corresponding check method

	// Check list type
	// Need to implement corresponding check method

	// Check set type
	// Need to implement corresponding check method

	// If none exists, return "none"
	return "none", nil
}

// Del deletes a key
func (m *DataStructureManager) Del(key string) (bool, error) {
	keyType, err := m.Type(key)
	if err != nil {
		return false, err
	}

	switch keyType {
	case "string":
		err := m.stringStruct.Del(key)
		return err == nil, err
	case "hash":
		return m.hashStruct.HDelAll(key)
	// Other type deletion methods
	// ...
	case "none":
		return false, nil
	default:
		return false, _const.ErrKeyNotFound
	}
}

// Expire sets key's expiration time
func (m *DataStructureManager) Expire(key string, ttl int64) (bool, error) {
	keyType, err := m.Type(key)
	if err != nil {
		return false, err
	}

	switch keyType {
	case "string":
		err := m.stringStruct.Expire(key, ttl)
		return err == nil, err
	case "hash":
		return m.hashStruct.HExpire(key, ttl)
	// Other type expiration settings
	// ...
	case "none":
		return false, nil
	default:
		return false, _const.ErrKeyNotFound
	}
}

// TTL gets key's remaining time to live
func (m *DataStructureManager) TTL(key string) (int64, error) {
	keyType, err := m.Type(key)
	if err != nil {
		return -2, err
	}

	switch keyType {
	case "string":
		return m.stringStruct.TTL(key)
	case "hash":
		return m.hashStruct.TTL(key)
	// Other type TTL retrieval
	// ...
	case "none":
		return -2, nil // Key doesn't exist
	default:
		return -2, _const.ErrKeyNotFound
	}
}

// Sync syncs data to disk
func (m *DataStructureManager) Sync() error {
	return m.db.Sync()
}

// Close closes the database
func (m *DataStructureManager) Close() error {
	return m.db.Close()
}

// Clean empties the database
func (m *DataStructureManager) Clean() {
	m.db.Clean()
}
