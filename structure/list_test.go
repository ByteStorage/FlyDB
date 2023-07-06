package structure

import (
	"os"
	"testing"

	"github.com/ByteStorage/FlyDB/config"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/ByteStorage/FlyDB/lib/randkv"
	"github.com/stretchr/testify/assert"
)

var listErr error

func initList() *ListStructure {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestListStructure")
	opts.DirPath = dir
	list, _ := NewListStructure(opts)
	return list
}

func TestListStructure_LPush(t *testing.T) {
	list := initList()
	defer list.db.Clean()

	// Test LPush function when the key exists
	listErr = list.LPush(randkv.GetTestKey(1), randkv.RandomValue(100))
	assert.Nil(t, listErr)

	// Test LPush function when the key does not exist
	listErr = list.LPush(randkv.GetTestKey(2), randkv.RandomValue(100))
	assert.Nil(t, listErr)
}

func TestListStructure_LPushs(t *testing.T) {
	list := initList()
	defer list.db.Clean()

	// Test LPushs function when the key exists
	listErr = list.LPushs(randkv.GetTestKey(1), randkv.RandomValue(100), randkv.RandomValue(100))
	assert.Nil(t, listErr)

	// Test LPushs function when the key does not exist
	listErr = list.LPushs(randkv.GetTestKey(2), randkv.RandomValue(100), randkv.RandomValue(100))
	assert.Nil(t, listErr)
}

func TestListStructure_RPush(t *testing.T) {
	list := initList()
	defer list.db.Clean()

	// Test RPush function when the key exists
	listErr = list.RPush(randkv.GetTestKey(1), randkv.RandomValue(100))
	assert.Nil(t, listErr)

	// Test RPush function when the key does not exist
	listErr = list.RPush(randkv.GetTestKey(2), randkv.RandomValue(100))
	assert.Nil(t, listErr)
}

func TestListStructure_RPushs(t *testing.T) {
	list := initList()
	defer list.db.Clean()

	// Test RPushs function when the key exists
	listErr = list.RPushs(randkv.GetTestKey(1), randkv.RandomValue(100), randkv.RandomValue(100))
	assert.Nil(t, listErr)

	// Test RPushs function when the key does not exist
	listErr = list.RPushs(randkv.GetTestKey(2), randkv.RandomValue(100), randkv.RandomValue(100))
	assert.Nil(t, listErr)
}

func TestListStructure_LPop(t *testing.T) {
	list := initList()
	defer list.db.Clean()

	// Test LPop function when the key exists
	listErr = list.LPush(randkv.GetTestKey(1), randkv.RandomValue(100))
	assert.Nil(t, listErr)
	value, err := list.LPop(randkv.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, value)

	// Test LPop function when the key does not exist
	_, err = list.LPop(randkv.GetTestKey(2))
	assert.Equal(t, err, _const.ErrKeyNotFound)

	// Test LPop function when the list is empty
	err = list.LPush(randkv.GetTestKey(3), randkv.RandomValue(100))
	assert.Nil(t, err)
	_, err = list.LPop(randkv.GetTestKey(3))
	assert.Nil(t, err)
	_, err = list.LPop(randkv.GetTestKey(3))
	assert.Equal(t, err, ErrListEmpty)
}

func TestListStructure_RPop(t *testing.T) {
	list := initList()
	defer list.db.Clean()

	// Test RPop function when the key exists
	listErr = list.RPush(randkv.GetTestKey(1), randkv.RandomValue(100))
	assert.Nil(t, listErr)
	value, err := list.RPop(randkv.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, value)

	// Test RPop function when the key does not exist
	_, err = list.RPop(randkv.GetTestKey(2))
	assert.Equal(t, err, _const.ErrKeyNotFound)

	// Test RPop function when the list is empty
	err = list.RPush(randkv.GetTestKey(3), randkv.RandomValue(100))
	assert.Nil(t, err)
	_, err = list.RPop(randkv.GetTestKey(3))
	assert.Nil(t, err)
	_, err = list.RPop(randkv.GetTestKey(3))
	assert.Equal(t, err, ErrListEmpty)
}

func TestListStructure_LRange(t *testing.T) {
	list := initList()
	defer list.db.Clean()

	// Test LRange function when the key exists
	listErr = list.LPush(randkv.GetTestKey(1), randkv.RandomValue(100))
	assert.Nil(t, listErr)
	values, err := list.LRange(randkv.GetTestKey(1), 0, 1)
	assert.Nil(t, err)
	assert.NotNil(t, values)

	// Test LRange function when the key does not exist
	_, err = list.LRange(randkv.GetTestKey(2), 0, 1)
	assert.Equal(t, err, _const.ErrKeyNotFound)

	// Test LRange function when the list is empty
	err = list.LPush(randkv.GetTestKey(3), randkv.RandomValue(100))
	assert.Nil(t, err)
	_, err = list.LPop(randkv.GetTestKey(3))
	assert.Nil(t, err)
	_, err = list.LRange(randkv.GetTestKey(3), 0, 1)
	assert.Equal(t, err, ErrListEmpty)
}

func TestListStructure_LLen(t *testing.T) {
	list := initList()
	defer list.db.Clean()

	// Test LLen function when the key exists
	listErr = list.LPush(randkv.GetTestKey(1), randkv.RandomValue(100))
	assert.Nil(t, listErr)
	len, err := list.LLen(randkv.GetTestKey(1))
	assert.Nil(t, err)
	assert.Equal(t, len, 1)

	// Test LLen function when the key does not exist
	_, err = list.LLen(randkv.GetTestKey(2))
	assert.Equal(t, err, _const.ErrKeyNotFound)
}

func TestListStructure_LRem(t *testing.T) {
	list := initList()
	defer list.db.Clean()

	// Test LRem function when the key exists
	listErr = list.LPush(randkv.GetTestKey(1), randkv.RandomValue(100))
	assert.Nil(t, listErr)
	listErr = list.LRem(randkv.GetTestKey(1), 1, randkv.RandomValue(100))
	assert.Nil(t, listErr)

	// Test LRem function when the key does not exist
	listErr = list.LRem(randkv.GetTestKey(2), 1, randkv.RandomValue(100))
	assert.Equal(t, listErr, _const.ErrKeyNotFound)
}

func TestListStructure_LSet(t *testing.T) {
	list := initList()
	defer list.db.Clean()

	// Test LSet function when the key exists
	listErr = list.LPush(randkv.GetTestKey(1), randkv.RandomValue(100))
	assert.Nil(t, listErr)
	listErr = list.LSet(randkv.GetTestKey(1), 0, randkv.RandomValue(200))
	assert.Nil(t, listErr)

	// Test LSet function when the key does not exist
	listErr = list.LSet(randkv.GetTestKey(2), 0, randkv.RandomValue(100))
	assert.Equal(t, listErr, _const.ErrKeyNotFound)

	// Test LSet function when the list is empty
	listErr = list.LPush(randkv.GetTestKey(3), randkv.RandomValue(100))
	assert.Nil(t, listErr)
	_, listErr = list.LPop(randkv.GetTestKey(3))
	assert.Nil(t, listErr)
	listErr = list.LSet(randkv.GetTestKey(3), 0, randkv.RandomValue(100))
	assert.Equal(t, listErr, ErrListEmpty)
}

func TestListStructure_LTrim(t *testing.T) {
	list := initList()
	defer list.db.Clean()

	// Test LTrim function when the key exists
	listErr = list.LPush(randkv.GetTestKey(1), randkv.RandomValue(100))
	assert.Nil(t, listErr)
	listErr = list.LTrim(randkv.GetTestKey(1), 0, 1)
	assert.Nil(t, listErr)

	// Test LTrim function when the key does not exist
	listErr = list.LTrim(randkv.GetTestKey(2), 0, 1)
	assert.Equal(t, listErr, _const.ErrKeyNotFound)

	// Test LTrim function when the list is empty
	listErr = list.LPush(randkv.GetTestKey(3), randkv.RandomValue(100))
	assert.Nil(t, listErr)
	_, listErr = list.LPop(randkv.GetTestKey(3))
	assert.Nil(t, listErr)
	listErr = list.LTrim(randkv.GetTestKey(3), 0, 1)
	assert.Equal(t, listErr, ErrListEmpty)
}

func TestListStructure_LIndex(t *testing.T) {
	list := initList()
	defer list.db.Clean()

	// Test LIndex function when the key exists
	listErr = list.LPush(randkv.GetTestKey(1), randkv.RandomValue(100))
	assert.Nil(t, listErr)
	value, err := list.LIndex(randkv.GetTestKey(1), 0)
	assert.Nil(t, err)
	assert.NotNil(t, value)

	// Test LIndex function when the key does not exist
	_, err = list.LIndex(randkv.GetTestKey(2), 0)
	assert.Equal(t, err, _const.ErrKeyNotFound)

	// Test LIndex function when the list is empty
	err = list.LPush(randkv.GetTestKey(3), randkv.RandomValue(100))
	assert.Nil(t, err)
	_, err = list.LPop(randkv.GetTestKey(3))
	assert.Nil(t, err)
	_, err = list.LIndex(randkv.GetTestKey(3), 0)
	assert.Equal(t, err, ErrListEmpty)
}

func TestListStructure_RPOPLPUSH(t *testing.T) {
	list := initList()
	defer list.db.Clean()

	// Test RPOPLPUSH function when the source list exists
	listErr = list.RPush(randkv.GetTestKey(1), randkv.RandomValue(100))
	assert.Nil(t, listErr)
	listErr = list.RPOPLPUSH(randkv.GetTestKey(1), randkv.GetTestKey(2))
	assert.Nil(t, listErr)

	// Test RPOPLPUSH function when the source list does not exist
	listErr = list.RPOPLPUSH(randkv.GetTestKey(3), randkv.GetTestKey(2))
	assert.Equal(t, listErr, _const.ErrKeyNotFound)

	// Test RPOPLPUSH function when the source list is empty
	listErr = list.RPush(randkv.GetTestKey(4), randkv.RandomValue(100))
	assert.Nil(t, listErr)
	_, listErr = list.RPop(randkv.GetTestKey(4))
	assert.Nil(t, listErr)
	listErr = list.RPOPLPUSH(randkv.GetTestKey(4), randkv.GetTestKey(2))
	assert.Equal(t, listErr, ErrListEmpty)
}

func TestListStructure_Integration(t *testing.T) {
	list := initList()
	defer list.db.Clean()

	// Create a key and use LPush to add some values
	key := randkv.GetTestKey(1)
	values := [][]byte{randkv.RandomValue(100), randkv.RandomValue(100), randkv.RandomValue(100)}
	for _, value := range values {
		listErr = list.RPush(key, value)
		assert.Nil(t, listErr)
	}

	// Use LLen to check the length of the list
	tmplen, err := list.LLen(key)
	assert.Nil(t, err)
	assert.Equal(t, tmplen, len(values))

	// Use LRange to get all values of the list and check if they are correct
	rangeValues, err := list.LRange(key, 0, -1)
	assert.Nil(t, err)
	assert.Equal(t, values, rangeValues)

	// Use LRem to remove a value and check if it is properly removed
	err = list.LRem(key, 1, values[0])
	assert.Nil(t, err)
	rangeValues, err = list.LRange(key, 0, -1)
	assert.Nil(t, err)
	assert.NotContains(t, rangeValues, values[0])

	// Use LSet to modify a value and check if it is properly modified
	newValue := randkv.RandomValue(100)
	err = list.LSet(key, 0, newValue)
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
	destination := randkv.GetTestKey(2)
	err = list.RPOPLPUSH(key, destination)
	assert.Nil(t, err)
	rangeValues, err = list.LRange(key, 0, -1)
	assert.Equal(t, ErrListEmpty, err)
	assert.Equal(t, len(rangeValues), 0)
	rangeValues, err = list.LRange(destination, 0, -1)
	assert.Nil(t, err)
	assert.Equal(t, len(rangeValues), 1)
}
