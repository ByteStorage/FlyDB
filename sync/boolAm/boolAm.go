package boolAm

import (
	"sync/atomic"
)

// The purpose of this package is to provide atomic operations on bool types,
// which are not directly supported by the built-in atomic package in Go.
// It achieves this by leveraging the uint32 type and
// modifying it to enable atomic operations on bool values.

// Boolean is a boolean value which can be accessed and modified atomically.
type Boolean uint32

// SetBoolAtomic sets the value of b.
// The operation is atomic and write the value atomically.
// If v is true, the value is set to 1.
// If v is false, the value is set to 0.
// The value of b is not guaranteed to be 1 or 0 if v is not true or false.
func (b *Boolean) SetBoolAtomic(v bool) {
	if v {
		atomic.StoreUint32((*uint32)(b), 1)
	} else {
		atomic.StoreUint32((*uint32)(b), 0)
	}
}

// GetBoolAtomic returns the value of b.
// The operation is atomic and read the value atomically.
func (b *Boolean) GetBoolAtomic() bool {
	return atomic.LoadUint32((*uint32)(b)) != 0
}
