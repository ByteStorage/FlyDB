package structure

import (
	"github.com/ByteStorage/FlyDB/db/engine"
)

// ListStructureInternal is the internal version of ListStructure using shared DB instance
type ListStructureInternal struct {
	db *engine.DB
}

// newListStructureInternal creates a new ListStructureInternal
func newListStructureInternal(db *engine.DB) *ListStructureInternal {
	return &ListStructureInternal{
		db: db,
	}
}

// SetStructureInternal is the internal version of SetStructure using shared DB instance
type SetStructureInternal struct {
	db *engine.DB
}

// newSetStructureInternal creates a new SetStructureInternal
func newSetStructureInternal(db *engine.DB) *SetStructureInternal {
	return &SetStructureInternal{
		db: db,
	}
}

// ZSetStructureInternal is the internal version of ZSetStructure using shared DB instance
type ZSetStructureInternal struct {
	db *engine.DB
}

// newZSetStructureInternal creates a new ZSetStructureInternal
func newZSetStructureInternal(db *engine.DB) *ZSetStructureInternal {
	return &ZSetStructureInternal{
		db: db,
	}
}

// BitmapStructureInternal is the internal version of BitmapStructure using shared DB instance
type BitmapStructureInternal struct {
	db *engine.DB
}

// newBitmapStructureInternal creates a new BitmapStructureInternal
func newBitmapStructureInternal(db *engine.DB) *BitmapStructureInternal {
	return &BitmapStructureInternal{
		db: db,
	}
}

// StreamStructureInternal is the internal version of StreamStructure using shared DB instance
type StreamStructureInternal struct {
	db *engine.DB
}

// newStreamStructureInternal creates a new StreamStructureInternal
func newStreamStructureInternal(db *engine.DB) *StreamStructureInternal {
	return &StreamStructureInternal{
		db: db,
	}
}

// ExpiringKeyInternal is the internal version of ExpiringKey using shared DB instance
type ExpiringKeyInternal struct {
	db *engine.DB
}

// newExpiringKeyInternal creates a new ExpiringKeyInternal
func newExpiringKeyInternal(db *engine.DB) *ExpiringKeyInternal {
	return &ExpiringKeyInternal{
		db: db,
	}
}
