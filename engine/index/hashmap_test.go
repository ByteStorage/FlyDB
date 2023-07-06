//package index
//
//import (
//	"testing"
//
//	"github.com/ByteStorage/FlyDB/engine/data"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestHashMap_Put(t *testing.T) {
//	hm := NewHashMap()
//
//	res1 := hm.Put(nil, &data.LogRecordPst{Fid: 1, Offset: 100})
//	assert.True(t, res1)
//
//	res2 := hm.Put([]byte("a"), &data.LogRecordPst{Fid: 1, Offset: 200})
//	assert.True(t, res2)
//}
//
//func TestHashMap_Get(t *testing.T) {
//	hm := NewHashMap()
//
//	res1 := hm.Put(nil, &data.LogRecordPst{Fid: 1, Offset: 100})
//	assert.True(t, res1)
//
//	pst1 := hm.Get(nil)
//	assert.Equal(t, uint32(1), pst1.Fid)
//	assert.Equal(t, int64(100), pst1.Offset)
//
//	res2 := hm.Put([]byte("a"), &data.LogRecordPst{Fid: 1, Offset: 200})
//	assert.True(t, res2)
//	res3 := hm.Put([]byte("a"), &data.LogRecordPst{Fid: 1, Offset: 300})
//	assert.True(t, res3)
//
//	pst2 := hm.Get([]byte("a"))
//	assert.Equal(t, uint32(1), pst2.Fid)
//	assert.Equal(t, int64(300), pst2.Offset)
//}
//
//func TestHashMap_Delete(t *testing.T) {
//	hm := NewHashMap()
//
//	res1 := hm.Put(nil, &data.LogRecordPst{Fid: 1, Offset: 100})
//	assert.True(t, res1)
//	res2 := hm.Delete(nil)
//	assert.True(t, res2)
//
//	res3 := hm.Put([]byte("abc"), &data.LogRecordPst{Fid: 11, Offset: 22})
//	assert.True(t, res3)
//	res4 := hm.Delete([]byte("abc"))
//	assert.True(t, res4)
//}
//
//func TestHashMap_Iterator(t *testing.T) {
//	hm1 := NewHashMap()
//	// 1. HashMap is empty
//
//	iter1 := hm1.Iterator(false)
//	assert.Equal(t, false, iter1.Valid())
//
//	// 2. HashMap is not empty
//	hm1.Put([]byte("abc"), &data.LogRecordPst{Fid: 1, Offset: 12})
//
//	iter2 := hm1.Iterator(false)
//	assert.True(t, iter2.Valid())
//	assert.NotNil(t, iter2.Key())
//	assert.NotNil(t, iter2.Value())
//	iter2.Next()
//	assert.Equal(t, false, iter2.Valid())
//
//	// 3. when there are multiple pieces of data
//	hm1.Put([]byte("bcd"), &data.LogRecordPst{Fid: 2, Offset: 12})
//	hm1.Put([]byte("efg"), &data.LogRecordPst{Fid: 3, Offset: 12})
//	hm1.Put([]byte("def"), &data.LogRecordPst{Fid: 4, Offset: 12})
//	iter3 := hm1.Iterator(false)
//	for iter3.Rewind(); iter3.Valid(); iter3.Next() {
//		assert.NotNil(t, iter3.Key())
//	}
//
//	iter4 := hm1.Iterator(true)
//	for iter4.Rewind(); iter4.Valid(); iter4.Next() {
//		assert.NotNil(t, iter4.Key())
//	}
//
//	// 4. Seek test
//	iter5 := hm1.Iterator(false)
//	for iter5.Seek([]byte("b")); iter5.Valid(); iter5.Next() {
//		assert.NotNil(t, iter5.Key())
//	}
//
//	iter6 := hm1.Iterator(true)
//	for iter6.Seek([]byte("d")); iter6.Valid(); iter6.Next() {
//		assert.NotNil(t, iter6.Key())
//	}
//
//}
