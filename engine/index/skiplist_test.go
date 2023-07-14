package index

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/engine/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSkipList_Put(t *testing.T) {
	sk := NewSkipList()

	res1 := sk.Put(nil, &data.LogRecordPst{Fid: 1, Offset: 100})
	assert.True(t, res1)

	res2 := sk.Put([]byte("a"), &data.LogRecordPst{Fid: 1, Offset: 200})
	assert.True(t, res2)
}

func TestSkipList_Get(t *testing.T) {
	sk := NewSkipList()

	res1 := sk.Put(nil, &data.LogRecordPst{Fid: 1, Offset: 100})
	assert.True(t, res1)

	pst1 := sk.Get(nil)
	assert.Equal(t, uint32(1), pst1.Fid)
	assert.Equal(t, int64(100), pst1.Offset)

	res2 := sk.Put([]byte("a"), &data.LogRecordPst{Fid: 1, Offset: 200})
	assert.True(t, res2)
	res3 := sk.Put([]byte("a"), &data.LogRecordPst{Fid: 1, Offset: 300})
	assert.True(t, res3)

	pst2 := sk.Get([]byte("a"))
	assert.Equal(t, uint32(1), pst2.Fid)
	assert.Equal(t, int64(300), pst2.Offset)
}

func TestSkipList_Delete(t *testing.T) {
	sk := NewSkipList()

	res1 := sk.Put(nil, &data.LogRecordPst{Fid: 1, Offset: 100})
	assert.True(t, res1)
	res2 := sk.Delete(nil)
	assert.True(t, res2)

	res3 := sk.Put([]byte("abc"), &data.LogRecordPst{Fid: 11, Offset: 22})
	assert.True(t, res3)
	res4 := sk.Delete([]byte("abc"))
	assert.True(t, res4)
	res5 := sk.Get([]byte("abc"))
	fmt.Println(res5 == nil)
	assert.Equal(t, (*data.LogRecordPst)(nil), res5)
}

func TestSkipList_Iterator(t *testing.T) {
	sk1 := NewSkipList()
	iter1 := sk1.Iterator(false)
	assert.Equal(t, false, iter1.Valid())

	sk1.Put([]byte("abc"), &data.LogRecordPst{Fid: 1, Offset: 12})
	iter2 := sk1.Iterator(false)
	assert.True(t, iter2.Valid())
	assert.NotNil(t, iter2.Key())
	assert.NotNil(t, iter2.Value())
	iter2.Next()
	assert.Equal(t, false, iter2.Valid())

	sk1.Put([]byte("bcd"), &data.LogRecordPst{Fid: 2, Offset: 12})
	sk1.Put([]byte("efg"), &data.LogRecordPst{Fid: 3, Offset: 12})
	sk1.Put([]byte("def"), &data.LogRecordPst{Fid: 4, Offset: 12})
	iter3 := sk1.Iterator(false)
	for iter3.Rewind(); iter3.Valid(); iter3.Next() {
		assert.NotNil(t, iter3.Key())
	}

	iter4 := sk1.Iterator(true)
	for iter4.Rewind(); iter4.Valid(); iter4.Next() {
		assert.NotNil(t, iter4.Key())
	}

	iter5 := sk1.Iterator(false)
	for iter5.Seek([]byte("b")); iter5.Valid(); iter5.Next() {
		assert.NotNil(t, iter5.Key())
	}

	iter6 := sk1.Iterator(true)
	for iter6.Seek([]byte("d")); iter6.Valid(); iter6.Next() {
		assert.NotNil(t, iter6.Key())
	}
}
