package index

import (
	"github.com/ByteStorage/FlyDB/engine/data"
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

func TestBTree_Iterator(t *testing.T) {
	bt1 := NewBTree()
	// 1. BTree 为空
	iter1 := bt1.Iterator(false)
	assert.Equal(t, false, iter1.Valid())

	// 2. BTree 不为空
	bt1.Put([]byte("abc"), &data.LogRecordPst{Fid: 1, Offset: 12})
	iter2 := bt1.Iterator(false)
	assert.True(t, iter2.Valid())
	assert.NotNil(t, iter2.Key())
	assert.NotNil(t, iter2.Value())
	iter2.Next()
	assert.Equal(t, false, iter2.Valid())

	// 3. 多条数据
	bt1.Put([]byte("bcd"), &data.LogRecordPst{Fid: 2, Offset: 12})
	bt1.Put([]byte("efg"), &data.LogRecordPst{Fid: 3, Offset: 12})
	bt1.Put([]byte("def"), &data.LogRecordPst{Fid: 4, Offset: 12})
	iter3 := bt1.Iterator(false)
	for iter3.Rewind(); iter3.Valid(); iter3.Next() {
		assert.NotNil(t, iter3.Key())
	}

	iter4 := bt1.Iterator(true)
	for iter4.Rewind(); iter4.Valid(); iter4.Next() {
		assert.NotNil(t, iter4.Key())
	}

	// 4. Seek test
	iter5 := bt1.Iterator(false)
	for iter5.Seek([]byte("b")); iter5.Valid(); iter5.Next() {
		assert.NotNil(t, iter5.Key())
	}

	iter6 := bt1.Iterator(true)
	for iter6.Seek([]byte("d")); iter6.Valid(); iter6.Next() {
		assert.NotNil(t, iter6.Key())
	}

}
