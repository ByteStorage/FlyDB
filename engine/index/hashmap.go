package index

//
//import (
//	"bytes"
//	"sort"
//	"sync"
//
//	"github.com/ByteStorage/FlyDB/engine/data"
//	"github.com/cornelk/hashmap"
//)
//
///*
//hashmap index, encapsulates the hashmap library of cornelk
//*/
//
//// HashMap struct
//type HashMap struct {
//	hashmap *hashmap.Map[string, *data.LogRecordPst]
//	// To ensure Thread safety
//	// multi thread writes need to be locked
//	lock *sync.RWMutex
//}
//
//// NewHashMap create a hashmap index
//func NewHashMap() *HashMap {
//	return &HashMap{
//		hashmap: hashmap.New[string, *data.LogRecordPst](),
//		lock:    new(sync.RWMutex),
//	}
//}
//
//// Implement the methods of the index interface
//// Put stores the data location information of key into the index
//func (hm *HashMap) Put(key []byte, pst *data.LogRecordPst) bool {
//	hm.lock.Lock()
//	defer hm.lock.Unlock()
//
//	hm.hashmap.Set(string(key), pst)
//	return true
//}
//
//// Get gains the data location of the key in the index
//func (hm *HashMap) Get(key []byte) *data.LogRecordPst {
//	value, ok := hm.hashmap.Get(string(key))
//	if !ok {
//		return nil
//	}
//	return value
//}
//
//// Delete deletes data location of one key in index
//func (hm *HashMap) Delete(key []byte) bool {
//	hm.lock.Lock()
//	hm.lock.Unlock()
//
//	return hm.hashmap.Del(string(key))
//}
//
//// Size returns the size of the data in index
//func (hm *HashMap) Size() int {
//	return hm.hashmap.Len()
//}
//
//// Iterator returns a index Iterator
//func (hm *HashMap) Iterator(reverse bool) Iterator {
//	// if the HashMap is empty, returns a default iterator
//	if hm.hashmap == nil {
//		return NewDefaultHashMapIterator(reverse)
//	}
//	hm.lock.RLock()
//	defer hm.lock.RUnlock()
//	// if the HashMap is not empty, returns a iterator
//	return NewHashMapIterator(hm.hashmap, reverse)
//}
//
//// HashMapIterator struct
//type HashMapIterator struct {
//	currIndex int     // The subscript position of the current traversal
//	reverse   bool    // Whether it is reverse traversal
//	values    []*Item // Key + Location index information
//}
//
//// create a default HashMap Iterator for the empty HashMap
//func NewDefaultHashMapIterator(reverse bool) *HashMapIterator {
//	return &HashMapIterator{
//		currIndex: 0,
//		reverse:   reverse,
//		values:    nil,
//	}
//}
//
//// create a HashMapIterator
//func NewHashMapIterator(hm *hashmap.Map[string, *data.LogRecordPst], reverse bool) *HashMapIterator {
//	// Use values slice to store all data in values
//	values := make([]*Item, hm.Len())
//
//	// count the number of elements in the values slice
//	var count int = 0
//	// store all data into an slice values
//	// We use range() method in the hashmap implement to do this
//	// define an operator method
//	saveFunc := func(key string, value *data.LogRecordPst) bool {
//		count++
//		item := &Item{
//			key: []byte(key),
//			pst: value,
//		}
//		values = append(values, item)
//		return true
//	}
//	// call range() method
//	hm.Range(saveFunc)
//
//	// filter out nil values
//	values = values[count:]
//
//	// if reverse needed, reverse the slice
//	if reverse {
//		for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
//			values[i], values[j] = values[j], values[i]
//		}
//	}
//
//	return &HashMapIterator{
//		currIndex: 0,
//		reverse:   reverse,
//		values:    values,
//	}
//}
//
//// Rewind goes back to the begining of the Iterator,ie. the index of the first data
//func (hmIt *HashMapIterator) Rewind() {
//	hmIt.currIndex = 0
//}
//
//// Seek finds a >= or <= target key according to the incoming key,
//// and starts traversing from this target key
//func (hmIt *HashMapIterator) Seek(key []byte) {
//	if hmIt.reverse {
//		hmIt.currIndex = sort.Search(len(hmIt.values), func(i int) bool {
//			return bytes.Compare(hmIt.values[i].key, key) <= 0
//		})
//	} else {
//		hmIt.currIndex = sort.Search(len(hmIt.values), func(i int) bool {
//			return bytes.Compare(hmIt.values[i].key, key) >= 0
//		})
//	}
//}
//
//// Next jumps to the next key
//func (hmIt *HashMapIterator) Next() {
//	hmIt.currIndex += 1
//}
//
//// Valid refers to whether it is valid, that is,
//// whether all keys have been traversed,
//// used to exit the traverse ==> true->yes false-->no
//func (hmIt *HashMapIterator) Valid() bool {
//	return hmIt.currIndex < len(hmIt.values)
//}
//
//// Key returns the key data at the current traversal position
//func (hmIt *HashMapIterator) Key() []byte {
//	return hmIt.values[hmIt.currIndex].key
//}
//
//// Value returns the value data of the current traversal position
//func (hmIt *HashMapIterator) Value() *data.LogRecordPst {
//	return hmIt.values[hmIt.currIndex].pst
//}
//
//// Close closes the iterator and releases the corresponding resources
//func (hmIt *HashMapIterator) Close() {
//	hmIt.values = nil
//}
