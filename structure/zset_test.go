package structure

import (
	"container/heap"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSortedSet(t *testing.T) {
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}
	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := ZSetNodes{}
	pq = make([]*ZSetNode, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &ZSetNode{
			Value:    value,
			Priority: priority,
			Index:    i,
		}
		i++
	}

	heap.Init(&pq)
	pq.Push(&ZSetNode{"Pineapple", 50, 0})
	//heap.Fix(&pq, len(pq)-1)
	//pq.update(pq[0], pq[0].value, 0)

	t.Log(pq)
}

func TestSortedSet_Bytes(t *testing.T) {
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}
	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := ZSetNodes{}
	pq = make([]*ZSetNode, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &ZSetNode{
			Value:    value,
			Priority: priority,
			Index:    i,
		}
		i++
	}

	heap.Init(&pq)
	b, err := pq.Bytes()
	assert.NoError(t, err)
	rb := ZSetNodes{}
	err = rb.FromBytes(b)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(rb, pq))
}
