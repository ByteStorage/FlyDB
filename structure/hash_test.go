package structure

import (
	"github.com/ByteStorage/FlyDB/config"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/ByteStorage/FlyDB/lib/randkv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func initHashDB() (*HashStructure, *config.Options) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestHashStructure")
	opts.DirPath = dir
	hash, _ := NewHashStructure(opts)
	return hash, &opts
}

func TestHashStructure_HGet(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok1)

	v1 := randkv.RandomValue(10)
	ok2, err := hash.HSet("1", []byte("field1"), v1)
	assert.Nil(t, err)
	assert.False(t, ok2)
	value1, err := hash.HGet("1", []byte("field1"))
	assert.Nil(t, err)
	assert.Equal(t, value1, v1)

	v2 := randkv.RandomValue(10)
	ok3, err := hash.HSet("1", []byte("field2"), v2)
	assert.Nil(t, err)
	assert.True(t, ok3)
	value2, err := hash.HGet("1", []byte("field2"))
	assert.Nil(t, err)
	assert.Equal(t, value2, v2)

	_, err = hash.HGet("1", []byte("field3"))
	assert.Equal(t, err, _const.ErrKeyNotFound)

	ok4, err := hash.HSet("2", "field1", "123123")
	assert.Nil(t, err)
	assert.True(t, ok4)

	v3, err := hash.HGet("2", "field1")
	assert.Nil(t, err)
	assert.Equal(t, v3, "123123")

}

func TestHashStructure_HMGet(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok1)

	v1 := randkv.RandomValue(10)
	ok2, err := hash.HSet("1", []byte("field1"), v1)
	assert.Nil(t, err)
	assert.False(t, ok2)

	v2 := randkv.RandomValue(10)
	ok3, err := hash.HSet("1", []byte("field2"), v2)
	assert.Nil(t, err)
	assert.True(t, ok3)

	mulVal, err := hash.HMGet("1", []byte("field1"), []byte("field2"))
	assert.Equal(t, v1, mulVal[0])
	assert.Equal(t, v2, mulVal[1])

	_, err = hash.HGet("1", []byte("field3"))
	assert.Equal(t, err, _const.ErrKeyNotFound)

}

func TestHashStructure_HDel(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok, err := hash.HDel("1", []byte("field1"))
	assert.Nil(t, err)
	assert.False(t, ok)

	ok1, err := hash.HSet("1", []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HDel("1", []byte("field1"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	v, err := hash.HGet("1", []byte("field1"))
	assert.Nil(t, err)
	assert.Nil(t, v)

	ok3, err := hash.HSet("1", []byte("field1"), "v111")
	assert.Nil(t, err)
	assert.True(t, ok3)

	ok4, err := hash.HSet("1", []byte("field2"), "v222")
	assert.Nil(t, err)
	assert.True(t, ok4)

}

func TestHashStructure_HExists(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), randkv.RandomValue(100))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HExists("1", []byte("field1"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HExists("1", []byte("field2"))
	assert.Nil(t, err)
	assert.False(t, ok3)

}

func TestHashStructure_HLen(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), randkv.RandomValue(100))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), randkv.RandomValue(100))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), randkv.RandomValue(100))
	assert.Nil(t, err)
	assert.True(t, ok3)

	l, err := hash.HLen("1")
	assert.Nil(t, err)
	assert.Equal(t, l, 3)
}

func TestHashStructure_HUpdate(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), randkv.RandomValue(100))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), randkv.RandomValue(100))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), randkv.RandomValue(100))
	assert.Nil(t, err)
	assert.True(t, ok3)

	ok4, err := hash.HUpdate("1", []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok4)

	ok5, err := hash.HUpdate("1", []byte("field2"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok5)

	ok6, err := hash.HUpdate("1", []byte("field4"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.False(t, ok6)

	l, err := hash.HLen("1")
	assert.Nil(t, err)
	assert.Equal(t, l, 3)
}

func TestHashStructure_HIncrBy(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok3)

	v1, err := hash.HIncrBy("1", []byte("field1"), 1)
	assert.Nil(t, err)
	assert.Equal(t, v1, int64(11))

	v2, err := hash.HIncrBy("1", []byte("field2"), -1)
	assert.Nil(t, err)
	assert.Equal(t, v2, int64(9))

	v3, err := hash.HIncrBy("1", []byte("field3"), 0)
	assert.Nil(t, err)
	assert.Equal(t, v3, int64(10))

	v4, err := hash.HIncrBy("1", []byte("field4"), 1)
	assert.Nil(t, err)
	assert.Equal(t, v4, int64(0))

}

func TestHashStructure_HIncrByFloat(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok3)

	v1, err := hash.HIncrByFloat("1", []byte("field1"), 1.1)
	assert.Nil(t, err)
	assert.Equal(t, v1, float64(11.1))

	v2, err := hash.HIncrByFloat("1", []byte("field2"), -1.1)
	assert.Nil(t, err)
	assert.Equal(t, v2, float64(8.9))

	v3, err := hash.HIncrByFloat("1", []byte("field3"), 0)
	assert.Nil(t, err)
	assert.Equal(t, v3, float64(10))

	v4, err := hash.HIncrByFloat("1", []byte("field4"), 1.1)
	assert.Nil(t, err)
	assert.Equal(t, v4, float64(0))

}

func TestHashStructure_HDecrBy(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok3)

	v1, err := hash.HDecrBy("1", []byte("field1"), 1)
	assert.Nil(t, err)
	assert.Equal(t, v1, int64(9))

	v2, err := hash.HDecrBy("1", []byte("field2"), 10)
	assert.Nil(t, err)
	assert.Equal(t, v2, int64(0))

	v3, err := hash.HDecrBy("1", []byte("field3"), 0)
	assert.Nil(t, err)
	assert.Equal(t, v3, int64(10))

	v4, err := hash.HDecrBy("1", []byte("field4"), 1)
	assert.Nil(t, err)
	assert.Equal(t, v4, int64(0))

}

func TestHashStructure_HStrLen(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), []byte("1000"))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), []byte("100"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok3)

	l1, err := hash.HStrLen("1", []byte("field1"))
	assert.Nil(t, err)
	assert.Equal(t, l1, 4)

	l2, err := hash.HStrLen("1", []byte("field2"))
	assert.Nil(t, err)
	assert.Equal(t, l2, 3)

	l3, err := hash.HStrLen("1", []byte("field3"))
	assert.Nil(t, err)
	assert.Equal(t, l3, 2)

	l4, err := hash.HStrLen("1", []byte("field4"))
	assert.Nil(t, err)
	assert.Equal(t, l4, 0)

}

func TestHashStructure_HMove(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), []byte("111-1000"))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), []byte("111-100"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), []byte("111-10"))
	assert.Nil(t, err)
	assert.True(t, ok3)

	ok4, err := hash.HSet("2", []byte("field1"), []byte("222-1000"))
	assert.Nil(t, err)
	assert.True(t, ok4)

	ok5, err := hash.HSet("2", []byte("field2"), []byte("222-100"))
	assert.Nil(t, err)
	assert.True(t, ok5)

	ok6, err := hash.HSet("2", []byte("field3"), []byte("222-10"))
	assert.Nil(t, err)
	assert.True(t, ok6)

	ok7, err := hash.HMove("2", "1", []byte("field1"))
	assert.Nil(t, err)
	assert.True(t, ok7)

	ok8, err := hash.HMove("2", "1", []byte("field2"))
	assert.Nil(t, err)
	assert.True(t, ok8)

	ok9, err := hash.HMove("2", "1", []byte("field3"))
	assert.Nil(t, err)
	assert.True(t, ok9)

	ok10, err := hash.HMove("2", "1", []byte("field4"))
	assert.Nil(t, err)
	assert.False(t, ok10)

	v1, err := hash.HGet("2", []byte("field1"))
	assert.Nil(t, err)
	assert.Equal(t, v1, []byte("111-1000"))

	v2, err := hash.HGet("2", []byte("field2"))
	assert.Nil(t, err)
	assert.Equal(t, v2, []byte("111-100"))

	v3, err := hash.HGet("2", []byte("field3"))
	assert.Nil(t, err)
	assert.Equal(t, v3, []byte("111-10"))

}

func TestHashStructure_HSetNX(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSetNX("1", []byte("field1"), []byte("1000"))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSetNX("1", []byte("field2"), []byte("100"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSetNX("1", []byte("field3"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok3)

	ok4, err := hash.HSetNX("1", []byte("field1"), []byte("1000"))
	assert.Nil(t, err)
	assert.False(t, ok4)

	ok5, err := hash.HSetNX("1", []byte("field2"), []byte("100"))
	assert.Nil(t, err)
	assert.False(t, ok5)

	ok6, err := hash.HSetNX("1", []byte("field3"), []byte("10"))
	assert.Nil(t, err)
	assert.False(t, ok6)

}

func TestHashStructure_HTypes(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), []byte("1000"))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), []byte("100"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok3)

	type1, err := hash.HTypes("1", []byte("field1"))
	assert.Nil(t, err)
	assert.Equal(t, type1, "hash")

	type2, err := hash.HTypes("1", []byte("field2"))
	assert.Nil(t, err)
	assert.Equal(t, type2, "hash")

	type3, err := hash.HTypes("1", []byte("field3"))
	assert.Nil(t, err)
	assert.Equal(t, type3, "hash")

	type4, err := hash.HTypes("1", []byte("field4"))
	assert.Equal(t, "", type4)
	assert.Equal(t, err, _const.ErrKeyNotFound)
}

func TestHashStructure_TTL(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", "field1", "123123")
	assert.Nil(t, err)
	assert.True(t, ok1)

	ttl, err := hash.TTL("1")
	assert.Nil(t, err)
	assert.Equal(t, ttl, int64(0))

	ok2, err := hash.HExpire("1", 2)
	assert.Nil(t, err)
	assert.True(t, ok2)

	time.Sleep(time.Second * 1)

	ttl, err = hash.TTL("1")
	assert.Nil(t, err)
	assert.Equal(t, ttl, int64(1))

	time.Sleep(time.Second * 3)

	ttl, err = hash.TTL("1")
	assert.NotNil(t, err)
	assert.Equal(t, ttl, int64(-1))

}

func TestHashStructure_Size(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", "field1", "11111")
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", "field2", "22222")
	assert.Nil(t, err)
	assert.True(t, ok2)

	value, err := hash.Size("1", "field1", "field2")
	assert.Nil(t, err)
	assert.Equal(t, value, "10B")

}

func TestHashStructure_Keys(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("qqqqqq", "!qqq!1", "11111")
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("qqqqqq", "!qqq!2", "22222")
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("qqqqqq", "!qqq!3", "33333")
	assert.Nil(t, err)
	assert.True(t, ok3)

	ok4, err := hash.HSet("qweqwe", "!qwe!1", "11111")
	assert.Nil(t, err)
	assert.True(t, ok4)

	ok5, err := hash.HSet("qweqwe", "!qwe!2", "22222")
	assert.Nil(t, err)
	assert.True(t, ok5)

	ok6, err := hash.HSet("qweqwe", "!qwe!3", "33333")
	assert.Nil(t, err)
	assert.True(t, ok6)

}

func TestHashStructure_GetFields(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("qqqqqq", "!qqq!1", "11111")
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("qqqqqq", "!qqq!2", "22222")
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("qqqqqq", "!qqq!3", "33333")
	assert.Nil(t, err)
	assert.True(t, ok3)

	ok11, err := hash.HSet("qweqwe", "!aaa!1", "11111")
	assert.Nil(t, err)
	assert.True(t, ok11)

	ok22, err := hash.HSet("qweqwe", "!aaa!2", "22222")
	assert.Nil(t, err)
	assert.True(t, ok22)

	ok33, err := hash.HSet("qweqwe", "!aaa!3", "33333")
	assert.Nil(t, err)
	assert.True(t, ok33)

	fields1 := hash.GetFields("qweqwe")
	assert.Equal(t, fields1, []string{"!aaa!1", "!aaa!2", "!aaa!3"})

	fields2 := hash.GetFields("qqqqqq")
	assert.Equal(t, fields2, []string{"!qqq!1", "!qqq!2", "!qqq!3"})

	keys, err := hash.Keys("*")
	assert.Nil(t, err)
	assert.Equal(t, keys, []string{"qqqqqq", "qweqwe"})

	keys, err = hash.Keys("q*")
	assert.Nil(t, err)
	assert.Equal(t, keys, []string{"qqqqqq", "qweqwe"})

	keys, err = hash.Keys("qq*")
	assert.Nil(t, err)
	assert.Equal(t, keys, []string{"qqqqqq"})

	keys, err = hash.Keys("q?*")
	assert.Nil(t, err)
	assert.Equal(t, keys, []string{"qqqqqq", "qweqwe"})

}

func TestHashStructure_HGetAllFieldAndValue(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("qqqqqq", "!qqq!1", "11111")
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("qqqqqq", "!qqq!2", "22222")
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("qqqqqq", "!qqq!3", "33333")
	assert.Nil(t, err)
	assert.True(t, ok3)

	ok11, err := hash.HSet("qweqwe", "!aaa!1", "11111")
	assert.Nil(t, err)
	assert.True(t, ok11)

	ok22, err := hash.HSet("qweqwe", "!aaa!2", "22222")
	assert.Nil(t, err)
	assert.True(t, ok22)

	ok33, err := hash.HSet("qweqwe", "!aaa!3", "33333")
	assert.Nil(t, err)
	assert.True(t, ok33)

	fv1, err := hash.HGetAllFieldAndValue("qqqqqq")
	assert.Nil(t, err)
	assert.Equal(t, fv1, map[string]interface{}{"!qqq!1": "11111", "!qqq!2": "22222", "!qqq!3": "33333"})

	fv2, err := hash.HGetAllFieldAndValue("qweqwe")
	assert.Nil(t, err)
	assert.Equal(t, fv2, map[string]interface{}{"!aaa!1": "11111", "!aaa!2": "22222", "!aaa!3": "33333"})

}

func TestHashStructure_HDelAll(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", "field1", "111111")
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", "field2", "222222")
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", "field3", "333333")
	assert.Nil(t, err)
	assert.True(t, ok3)

	//ok4, err := hash.HDelAll("1")
	//fmt.Println(ok4, err)

	//v1, err := hash.HGet("1", "field1")
	//fmt.Println(v1, err)
	//v2, err := hash.HGet("1", "field2")
	//fmt.Println(v2, err)
	//v3, err := hash.HGet("1", "field3")
	//fmt.Println(v3, err)

	//ttl, err := hash.TTL("1")
	//fmt.Println("--", ttl, err)
}
