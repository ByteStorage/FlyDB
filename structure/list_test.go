package structure

import (
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/ByteStorage/FlyDB/lib/randkv"
	"os"
	"testing"
	"time"

	"github.com/ByteStorage/FlyDB/config"
	"github.com/stretchr/testify/assert"
)

var listErr error

func initList() (*ListStructure, *config.Options) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestListStructure")
	opts.DirPath = dir
	list, _ := NewListStructure(opts)
	return list, &opts
}

func TestListStructure_TTL(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	err := list.LPush("2", "123123", 2)
	assert.Nil(t, err)

	ttl, err := list.TTL("2")
	assert.Nil(t, err)
	assert.Equal(t, ttl, int64(2))

	time.Sleep(1 * time.Second)
	ttl, err = list.TTL("2")
	assert.Nil(t, err)
	assert.Equal(t, ttl, int64(1))

	time.Sleep(2 * time.Second)
	ttl, err = list.TTL("1")
	assert.NotNil(t, err)
	assert.Equal(t, ttl, int64(-1))

	err = list.LPush("1", "123123", 0)
	assert.Nil(t, err)

	ttl1, err := list.TTL("1")
	assert.Nil(t, err)
	assert.Equal(t, ttl1, int64(0))
}

func TestListStructure_Size(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	err = list.LPush("2", "123123", 0)
	assert.Nil(t, err)

	err = list.LPush("2", "1233", 0)
	assert.Nil(t, err)

	value, err := list.Size("2")
	assert.Nil(t, err)
	assert.Equal(t, value, "10B")
}

func TestListStructure_LPush(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	// Test LPush function when the key exists
	listErr = list.LPush("key", "w", 10)
	assert.Nil(t, listErr)

	//// Test LPush function when the key does not exist
	//listErr = list.LPush(string(randkv.GetTestKey(2)), randkv.RandomValue(100))
	//assert.Nil(t, listErr)
}

func TestListStructure_LPushs(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	// Test LPushs function when the key exists
	listErr = list.LPushs(string(randkv.GetTestKey(1)), 0, randkv.RandomValue(100), randkv.RandomValue(100))
	assert.Nil(t, listErr)

	// Test LPushs function when the key does not exist
	listErr = list.LPushs(string(randkv.GetTestKey(2)), 0, randkv.RandomValue(100), randkv.RandomValue(100))
	assert.Nil(t, listErr)
}

func TestListStructure_RPush(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	// Test RPush function when the key exists
	listErr = list.RPush(string(randkv.GetTestKey(1)), randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)

	// Test RPush function when the key does not exist
	listErr = list.RPush(string(randkv.GetTestKey(2)), randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)
}

func TestListStructure_RPushs(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	// Test RPushs function when the key exists
	listErr = list.RPushs(string(randkv.GetTestKey(1)), 0, randkv.RandomValue(100), randkv.RandomValue(100))
	assert.Nil(t, listErr)

	// Test RPushs function when the key does not exist
	listErr = list.RPushs(string(randkv.GetTestKey(2)), 0, randkv.RandomValue(100), randkv.RandomValue(100))
	assert.Nil(t, listErr)
}

func TestListStructure_LPop(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	// Test LPop function when the key exists
	listErr = list.LPush(string(randkv.GetTestKey(1)), randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)
	listErr = list.LPush(string(randkv.GetTestKey(1)), "www", 0)
	value, err := list.LPop(string(randkv.GetTestKey(1)))
	assert.Nil(t, err)
	assert.NotNil(t, value)

	// Test LPop function when the key does not exist
	_, err = list.LPop(string(randkv.GetTestKey(2)))
	assert.Equal(t, err, _const.ErrKeyNotFound)

	// Test LPop function when the list is empty
	err = list.LPush(string(randkv.GetTestKey(3)), randkv.RandomValue(100), 0)
	assert.Nil(t, err)
	_, err = list.LPop(string(randkv.GetTestKey(3)))
	assert.Nil(t, err)
	_, err = list.LPop(string(randkv.GetTestKey(3)))
	assert.Equal(t, err, ErrListEmpty)
}

func TestListStructure_RPop(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	// Test RPop function when the key exists
	listErr = list.RPush(string(randkv.GetTestKey(1)), randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)
	value, err := list.RPop(string(randkv.GetTestKey(1)))
	assert.Nil(t, err)
	assert.NotNil(t, value)

	// Test RPop function when the key does not exist
	_, err = list.RPop(string(randkv.GetTestKey(2)))
	assert.Equal(t, err, _const.ErrKeyNotFound)

	// Test RPop function when the list is empty
	err = list.RPush(string(randkv.GetTestKey(3)), randkv.RandomValue(100), 0)
	assert.Nil(t, err)
	_, err = list.RPop(string(randkv.GetTestKey(3)))
	assert.Nil(t, err)
	_, err = list.RPop(string(randkv.GetTestKey(3)))
	assert.Equal(t, err, ErrListEmpty)
}

func TestListStructure_LRange(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	// Test LRange function when the key exists
	listErr = list.LPush(string(randkv.GetTestKey(1)), randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)
	values, err := list.LRange(string(randkv.GetTestKey(1)), 0, 1)
	assert.Nil(t, err)
	assert.NotNil(t, values)

	// Test LRange function when the key does not exist
	_, err = list.LRange(string(randkv.GetTestKey(2)), 0, 1)
	assert.Equal(t, err, _const.ErrKeyNotFound)

	// Test LRange function when the list is empty
	err = list.LPush(string(randkv.GetTestKey(3)), randkv.RandomValue(100), 0)
	assert.Nil(t, err)
	_, err = list.LPop(string(randkv.GetTestKey(3)))
	assert.Nil(t, err)
	_, err = list.LRange(string(randkv.GetTestKey(3)), 0, 1)
	assert.Equal(t, err, ErrListEmpty)
}

func TestListStructure_LLen(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	// Test LLen function when the key exists
	listErr = list.LPush(string(randkv.GetTestKey(1)), randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)
	len, err := list.LLen(string(randkv.GetTestKey(1)))
	assert.Nil(t, err)
	assert.Equal(t, len, 1)

	// Test LLen function when the key does not exist
	_, err = list.LLen(string(randkv.GetTestKey(2)))
	assert.Equal(t, err, _const.ErrKeyNotFound)
}

func TestListStructure_LRem(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	// Test LRem function when the key exists
	listErr = list.LPush(string(randkv.GetTestKey(1)), randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)
	listErr = list.LRem(string(randkv.GetTestKey(1)), 1, randkv.RandomValue(100))
	assert.Nil(t, listErr)

	// Test LRem function when the key does not exist
	listErr = list.LRem(string(randkv.GetTestKey(2)), 1, randkv.RandomValue(100))
	assert.Equal(t, listErr, _const.ErrKeyNotFound)
}

func TestListStructure_LSet(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	// Test LSet function when the key exists
	listErr = list.LPush(string(randkv.GetTestKey(1)), randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)
	listErr = list.LSet(string(randkv.GetTestKey(1)), 0, randkv.RandomValue(200), 0)
	assert.Nil(t, listErr)

	// Test LSet function when the key does not exist
	listErr = list.LSet(string(randkv.GetTestKey(2)), 0, randkv.RandomValue(100), 0)
	assert.Equal(t, listErr, _const.ErrKeyNotFound)

	// Test LSet function when the list is empty
	listErr = list.LPush(string(randkv.GetTestKey(3)), randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)
	_, listErr = list.LPop(string(randkv.GetTestKey(3)))
	assert.Nil(t, listErr)
	listErr = list.LSet(string(randkv.GetTestKey(3)), 0, randkv.RandomValue(100), 0)
	assert.Equal(t, listErr, ErrIndexOutOfRange)
}

func TestListStructure_LTrim(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	// Test LTrim function when the key exists
	listErr = list.LPush(string(randkv.GetTestKey(1)), randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)
	listErr = list.LTrim(string(randkv.GetTestKey(1)), 0, 1)
	assert.Nil(t, listErr)

	// Test LTrim function when the key does not exist
	listErr = list.LTrim(string(randkv.GetTestKey(2)), 0, 1)
	assert.Equal(t, listErr, _const.ErrKeyNotFound)

	// Test LTrim function when the list is empty
	listErr = list.LPush(string(randkv.GetTestKey(3)), randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)
	_, listErr = list.LPop(string(randkv.GetTestKey(3)))
	assert.Nil(t, listErr)
	listErr = list.LTrim(string(randkv.GetTestKey(3)), 0, 1)
	assert.Equal(t, listErr, ErrListEmpty)
}

func TestListStructure_LIndex(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	// Test LIndex function when the key exists
	listErr = list.LPush(string(randkv.GetTestKey(1)), randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)
	value, err := list.LIndex(string(randkv.GetTestKey(1)), 0)
	assert.Nil(t, err)
	assert.NotNil(t, value)

	// Test LIndex function when the key does not exist
	_, err = list.LIndex(string(randkv.GetTestKey(2)), 0)
	assert.Equal(t, err, _const.ErrKeyNotFound)

	// Test LIndex function when the list is empty
	err = list.LPush(string(randkv.GetTestKey(3)), randkv.RandomValue(100), 0)
	assert.Nil(t, err)
	_, err = list.LPop(string(randkv.GetTestKey(3)))
	assert.Nil(t, err)
	_, err = list.LIndex(string(randkv.GetTestKey(3)), 0)
	assert.Equal(t, err, ErrListEmpty)
}

func TestListStructure_RPOPLPUSH(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	// Test RPOPLPUSH function when the source list exists
	listErr = list.RPush(string(randkv.GetTestKey(1)), randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)
	listErr = list.RPOPLPUSH(string(randkv.GetTestKey(1)), string(randkv.GetTestKey(2)), 0)
	assert.Nil(t, listErr)

	// Test RPOPLPUSH function when the source list does not exist
	listErr = list.RPOPLPUSH(string(randkv.GetTestKey(3)), string(randkv.GetTestKey(2)), 0)
	assert.Equal(t, listErr, _const.ErrKeyNotFound)

	// Test RPOPLPUSH function when the source list is empty
	listErr = list.RPush(string(randkv.GetTestKey(4)), randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)
	_, listErr = list.RPop(string(randkv.GetTestKey(4)))
	assert.Nil(t, listErr)
	listErr = list.RPOPLPUSH(string(randkv.GetTestKey(4)), string(randkv.GetTestKey(2)), 0)
	assert.Equal(t, listErr, ErrListEmpty)
}

func TestListStructure_Integration(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	// Create a key and use LPush to add some values
	key := string(randkv.GetTestKey(1))
	values := [][]byte{randkv.RandomValue(100), randkv.RandomValue(100), randkv.RandomValue(100)}
	for _, value := range values {
		listErr = list.RPush(key, value, 0)
		assert.Nil(t, listErr)
	}

	// Use LLen to check the length of the list
	tmplen, err := list.LLen(key)
	assert.Nil(t, err)
	assert.Equal(t, tmplen, len(values))

	// Use LRange to get all values of the list and check if they are correct
	rangeValues, err := list.LRange(key, 0, -1)
	assert.Nil(t, err)
	bytesRangeValues := make([][]byte, len(rangeValues))
	for i := 0; i < len(rangeValues); i++ {
		bytesRangeValues[i] = rangeValues[i].([]byte)
	}

	assert.Equal(t, values, bytesRangeValues)

	// Use LRem to remove a value and check if it is properly removed
	err = list.LRem(key, 1, values[0])
	assert.Nil(t, err)
	rangeValues, err = list.LRange(key, 0, -1)
	assert.Nil(t, err)
	assert.NotContains(t, rangeValues, values[0])

	// Use LSet to modify a value and check if it is properly modified
	newValue := randkv.RandomValue(100)
	err = list.LSet(key, 0, newValue, 0)
	assert.Nil(t, err)
	rangeValues, err = list.LRange(key, 0, -1)
	assert.Nil(t, err)
	assert.Contains(t, rangeValues, newValue)

	// Use LTrim to trim the list and check if it is properly trimmed
	err = list.LTrim(key, 0, 0)
	assert.Nil(t, err)
	rangeValues, err = list.LRange(key, 0, -1)
	assert.Nil(t, err)
	assert.Equal(t, len(rangeValues), 1)

	// Use RPOPLPUSH to move a value to another list and check if it is properly moved
	destination := string(randkv.GetTestKey(2))
	err = list.RPOPLPUSH(key, destination, 0)
	assert.Nil(t, err)
	rangeValues, err = list.LRange(key, 0, -1)
	assert.Equal(t, ErrListEmpty, err)
	assert.Equal(t, len(rangeValues), 0)
	rangeValues, err = list.LRange(destination, 0, -1)
	assert.Nil(t, err)
	assert.Equal(t, len(rangeValues), 1)
}

func TestListStructure_Keys(t *testing.T) {
	list, _ := initList()
	defer list.db.Clean()

	listErr = list.LPush("1", randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)

	listErr = list.LPush("2", randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)

	listErr = list.LPush("q", randkv.RandomValue(100), 0)
	assert.Nil(t, listErr)

	keys, err := list.Keys()
	assert.Nil(t, err)
	assert.Equal(t, len(keys), 3)
}
