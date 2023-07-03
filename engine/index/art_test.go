package index

import (
	"github.com/ByteStorage/FlyDB/engine/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdaptiveRadixTree_Put(t *testing.T) {
	art := NewART()
	art.Put([]byte("key-1"), &data.LogRecordPst{Fid: 1, Offset: 12})
	art.Put([]byte("key-1"), &data.LogRecordPst{Fid: 1, Offset: 12})
	art.Put([]byte("key-1"), &data.LogRecordPst{Fid: 1, Offset: 12})
}

func TestAdaptiveRadixTree_Get(t *testing.T) {
	art := NewART()
	art.Put([]byte("key-1"), &data.LogRecordPst{Fid: 1, Offset: 12})
	pst := art.Get([]byte("key-1"))
	assert.NotNil(t, pst)

	pst1 := art.Get([]byte("key-not-exist"))
	assert.Nil(t, pst1)

	art.Put([]byte("key-1"), &data.LogRecordPst{Fid: 12, Offset: 123})
	pst2 := art.Get([]byte("key-1"))
	assert.NotNil(t, pst2)
}

func TestAdaptiveRadixTree_Delete(t *testing.T) {
	art := NewART()

	res := art.Delete([]byte("key-not-exist"))
	assert.False(t, res)

	art.Put([]byte("key-1"), &data.LogRecordPst{Fid: 1, Offset: 12})
	res1 := art.Delete([]byte("key-1"))
	assert.True(t, res1)

	pst := art.Get([]byte("key-1"))
	assert.Nil(t, pst)
}

func TestAdaptiveRadixTree_Size(t *testing.T) {
	art := NewART()

	assert.Equal(t, 0, art.Size())

	art.Put([]byte("key-1"), &data.LogRecordPst{Fid: 1, Offset: 12})
	art.Put([]byte("key-2"), &data.LogRecordPst{Fid: 1, Offset: 12})
	art.Put([]byte("key-1"), &data.LogRecordPst{Fid: 12, Offset: 123})
	assert.Equal(t, 2, art.Size())
}

func TestAdaptiveRadixTree_Iterator(t *testing.T) {
	art := NewART()

	art.Put([]byte("cdef"), &data.LogRecordPst{Fid: 1, Offset: 12})
	art.Put([]byte("bcde"), &data.LogRecordPst{Fid: 1, Offset: 12})
	art.Put([]byte("abcd"), &data.LogRecordPst{Fid: 1, Offset: 12})
	art.Put([]byte("bdfg"), &data.LogRecordPst{Fid: 1, Offset: 12})

	iter := art.Iterator(false)
	for iter.Rewind(); iter.Valid(); iter.Next() {
		assert.NotNil(t, iter.Key())
		assert.NotNil(t, iter.Value())
	}

}
