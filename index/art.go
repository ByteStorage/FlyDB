package index

import (
	art "github.com/plar/go-adaptive-radix-tree"
	"github.com/qishenonly/flydb/data"
	"sync"
)

// Adaptive Radix Tree Index
// https://github.com/plar/go-adaptive-radix-tree
type AdaptiveRadixTree struct {
	tree art.Tree
	lock *sync.RWMutex
}

// NewART Initializes the adaptive radix tree index
func NewART() *AdaptiveRadixTree {
	return &AdaptiveRadixTree{
		tree: art.New(),
		lock: new(sync.RWMutex),
	}
}

func (artree *AdaptiveRadixTree) Put(key []byte, pst *data.LogRecordPst) bool {}

func (artree *AdaptiveRadixTree) Get(key []byte) *data.LogRecordPst {}

func (artree *AdaptiveRadixTree) Delete(key []byte) bool {}

func (artree *AdaptiveRadixTree) Size() int {}

func (artree *AdaptiveRadixTree) Iterator(reverse bool) Iterator {}
