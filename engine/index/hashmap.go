package index

import (
	"bytes"
	"sort"
	"sync"

	"github.com/ByteStorage/FlyDB/engine/data"
	"github.com/cornelk/hashmap"
)

/*
hashmap index, encapsulates the hashmap library of cornelk
*/

// HashMap struct
type HashMap struct {
	hashmap *hashmap.Map[string, *data.LogRecordPst]
	// To ensure Thread safety
	// multi thread writes need to be locked
	lock *sync.RWMutex
}

// NewHashMap create a hashmap index
func NewHashMap() *HashMap {
	return &HashMap{
		hashmap: hashmap.New[string, *data.LogRecordPst](),
		lock:    new(sync.RWMutex),
	}
}

// Implement the methods of the index interface
// Put 向索引中存储 key 对应的数据位置信息
func (hm *HashMap) Put(key []byte, pst *data.LogRecordPst) bool {
	hm.lock.Lock()
	defer hm.lock.Unlock()

	hm.hashmap.Set(string(key), pst)
	return true
}

// Get 根据 key 取出对应的索引位置信息
func (hm *HashMap) Get(key []byte) *data.LogRecordPst {
	value, ok := hm.hashmap.Get(string(key))
	if !ok {
		return nil
	}
	return value
}

// Delete 根据 key 删除对应的索引位置信息
func (hm *HashMap) Delete(key []byte) bool {
	hm.lock.Lock()
	hm.lock.Unlock()

	return hm.hashmap.Del(string(key))
}

// Size 索引中的数据量
func (hm *HashMap) Size() int {
	return hm.hashmap.Len()
}

// Iterator 索引迭代器
func (hm *HashMap) Iterator(reverse bool) Iterator {
	if hm.hashmap == nil {
		return nil
	}
	hm.lock.RLock()
	defer hm.lock.RUnlock()
	return NewHashMapIterator(hm, reverse)
}

// HashMapIterator struct
type HashMapIterator struct {
	currIndex int     // The subscript position of the current traversal
	reverse   bool    // Whether it is reverse traversal
	values    []*Item // Key + Location index information
}

// create a HashMapIterator
func NewHashMapIterator(hm *HashMap, reverse bool) *HashMapIterator {
	values := make([]*Item, hm.Size())

	// store all data into an slice values
	// 使用hashmap实现中的range函数来做
	// 需要定义一个操作函数
	saveFunc := func(key string, value *data.LogRecordPst) bool {
		item := &Item{
			key: []byte(key),
			pst: value,
		}
		values = append(values, item)
		return true
	}
	// call range() method
	hm.hashmap.Range(saveFunc)
	// if reverse needed, reverse the slice
	if reverse {
		for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
			values[i], values[j] = values[j], values[i]
		}
	}

	return &HashMapIterator{
		currIndex: 0,
		reverse:   reverse,
		values:    values,
	}
}

// Rewind 重新回到迭代器的起点，即第一个数据
func (hmIt *HashMapIterator) Rewind() {
	hmIt.currIndex = 0
}

// Seek 根据传入的 key 查找到一个 >= 或 <= 的目标 key，从这个目标 key 开始遍历
func (hmIt *HashMapIterator) Seek(key []byte) {
	if hmIt.reverse {
		hmIt.currIndex = sort.Search(len(hmIt.values), func(i int) bool {
			return bytes.Compare(hmIt.values[i].key, key) <= 0
		})
	} else {
		hmIt.currIndex = sort.Search(len(hmIt.values), func(i int) bool {
			return bytes.Compare(hmIt.values[i].key, key) >= 0
		})
	}
}

// Next 跳转到下一个 key
func (hmIt *HashMapIterator) Next() {
	hmIt.currIndex += 1
}

// Valid 是否有效，即是否已经遍历完了所有 key，用于退出遍历 ==> true->是  false-->否
func (hmIt *HashMapIterator) Valid() bool {
	return hmIt.currIndex < len(hmIt.values)
}

// Key 当前遍历位置的 key 数据
func (hmIt *HashMapIterator) Key() []byte {
	return hmIt.values[hmIt.currIndex].key
}

// Value 当前遍历位置的 value 数据
func (hmIt *HashMapIterator) Value() *data.LogRecordPst {
	return hmIt.values[hmIt.currIndex].pst
}

// Close 关闭迭代器，释放相应资源
func (hmIt *HashMapIterator) Close() {
	hmIt.values = nil
}
