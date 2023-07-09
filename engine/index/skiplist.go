package index

import (
	"bytes"
	"github.com/ByteStorage/FlyDB/engine/data"
	"github.com/chen3feng/stl4go"
	"sort"
	"sync"
)

// SkipList Memory Index
// based on https://github.com/chen3feng/stl4go
type SkipList struct {
	list *stl4go.SkipList[[]byte, *data.LogRecordPst]
	lock *sync.RWMutex
}

// NewSkipList Initialize the SkipList index
func NewSkipList() *SkipList {
	return &SkipList{
		list: stl4go.NewSkipListFunc[[]byte, *data.LogRecordPst](Compare),
		lock: new(sync.RWMutex),
	}
}

// Put Inserts a key-value pair into the SkipList index
func (sl *SkipList) Put(key []byte, pst *data.LogRecordPst) bool {
	sl.lock.Lock()
	defer sl.lock.Unlock()

	sl.list.Insert(key, pst)
	return true
}

// Get Gets the value corresponding to the key from the SkipList index
func (sl *SkipList) Get(key []byte) *data.LogRecordPst {
	sl.lock.RLock()
	defer sl.lock.RUnlock()

	res := sl.list.Find(key)
	if res != nil {
		return *res
	}
	return nil
}

// Delete Deletes the key-value pair corresponding to the key from the SkipList index
func (sl *SkipList) Delete(key []byte) bool {
	sl.lock.Lock()
	defer sl.lock.Unlock()

	return sl.list.Remove(key)
}

// Size Gets the number of key-value pairs in the SkipList index
func (sl *SkipList) Size() int {
	return sl.list.Len()
}

// Iterator Gets the iterator of the SkipList index
// If the reverse is true, the iterator is traversed in reverse order,
// otherwise it is traversed in order
func (sl *SkipList) Iterator(reverse bool) Iterator {
	sl.lock.RLock()
	defer sl.lock.RUnlock()
	return NewSkipListIterator(sl, reverse)
}

type SkipListIterator struct {
	currIndex int
	reverse   bool
	values    []*Item
}

// NewSkipListIterator Initializes the SkipList index iterator
func NewSkipListIterator(sl *SkipList, reverse bool) *SkipListIterator {
	// Estimate the expected slice capacity based on skip list size
	expectedSize := sl.Size()

	// Initialize with empty slice and expected capacity
	values := make([]*Item, 0, expectedSize)

	// for each operation
	saveToValues := func(K []byte, V *data.LogRecordPst) {
		item := &Item{
			key: K,
			pst: V,
		}
		values = append(values, item)
	}
	sl.list.ForEach(saveToValues)

	// Reverse the values slice if reverse is true
	if reverse {
		for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
			values[i], values[j] = values[j], values[i]
		}
	}

	return &SkipListIterator{
		currIndex: 0,
		reverse:   reverse,
		values:    values,
	}
}

// Rewind Resets the iterator to the beginning
func (sl *SkipListIterator) Rewind() {
	sl.currIndex = 0
}

// Seek Positions the iterator to the first key
// that is greater or equal to the specified key
func (sl *SkipListIterator) Seek(key []byte) {
	// binary search
	if sl.reverse {
		sl.currIndex = sort.Search(len(sl.values), func(i int) bool {
			return bytes.Compare(sl.values[i].key, key) <= 0
		})
	} else {
		sl.currIndex = sort.Search(len(sl.values), func(i int) bool {
			return bytes.Compare(sl.values[i].key, key) >= 0
		})
	}
}

// Next Positions the iterator to the next key
// If the iterator is positioned at the last key,
// the iterator is positioned to the start of the iterator
func (sl *SkipListIterator) Next() {
	sl.currIndex += 1
}

// Valid Determines whether the iterator is positioned at a valid key
func (sl *SkipListIterator) Valid() bool {
	return sl.currIndex < len(sl.values)
}

// Key Gets the key at the current iterator position
func (sl *SkipListIterator) Key() []byte {
	return sl.values[sl.currIndex].key
}

// Value Gets the value at the current iterator position
func (sl *SkipListIterator) Value() *data.LogRecordPst {
	return sl.values[sl.currIndex].pst
}

// Close Closes the iterator
func (sl *SkipListIterator) Close() {
	sl.values = nil
}
