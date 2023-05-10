package index

import (
	art "github.com/plar/go-adaptive-radix-tree"
	"sync"
)

// Adaptive Radix Tree Index
// https://github.com/plar/go-adaptive-radix-tree
type AdaptiveRadixTree struct {
	tree art.Tree
	lock *sync.RWMutex
}
