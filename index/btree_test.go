package index

import (
	"flydb/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBTree_Put(t *testing.T) {
	bt := NewBTree()

	res1 := bt.Put(nil, &data.LogRecordPst{Fid: 1, Offset: 100})
	assert.True(t, res1)

	res2 := bt.Put([]byte("a"), &data.LogRecordPst{Fid: 1, Offset: 200})
	assert.True(t, res2)
}

func TestBTree_Get(t *testing.T) {
	bt := NewBTree()

	res1 := bt.Put(nil, &data.LogRecordPst{Fid: 1, Offset: 100})
	assert.True(t, res1)

	pst1 := bt.Get(nil)
	assert.Equal(t, uint32(1), pst1.Fid)
	assert.Equal(t, int64(100), pst1.Offset)

	res2 := bt.Put([]byte("a"), &data.LogRecordPst{Fid: 1, Offset: 200})
	assert.True(t, res2)
	res3 := bt.Put([]byte("a"), &data.LogRecordPst{Fid: 1, Offset: 300})
	assert.True(t, res3)

	pst2 := bt.Get([]byte("a"))
	assert.Equal(t, uint32(1), pst2.Fid)
	assert.Equal(t, int64(300), pst2.Offset)
}

func TestBTree_Delete(t *testing.T) {
	bt := NewBTree()

	res1 := bt.Put(nil, &data.LogRecordPst{Fid: 1, Offset: 100})
	assert.True(t, res1)
	res2 := bt.Delete(nil)
	assert.True(t, res2)

	res3 := bt.Put([]byte("abc"), &data.LogRecordPst{Fid: 11, Offset: 22})
	assert.True(t, res3)
	res4 := bt.Delete([]byte("abc"))
	assert.True(t, res4)
}
