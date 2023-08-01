package structure

import (
	"os"
	"testing"
	"time"

	"github.com/ByteStorage/FlyDB/config"
	"github.com/stretchr/testify/assert"
)

var expireErr error

func initExpire() (*ExpireStructure, *config.Options) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestExpireStructure")
	opts.DirPath = dir
	expire, _ := NewExpireStructure(opts)
	return expire, &opts
}

func TestExpireStructure_EXPIRE(t *testing.T) {
	expire, _ := initExpire()
	defer expire.db.Close()

	key := "test_key"
	seconds := int64(2)

	// Test EXPIRE function
	expireErr = expire.EXPIRE(key, seconds)
	assert.Nil(t, expireErr)

	// Wait for the key to expire
	time.Sleep(time.Second * time.Duration(seconds+1))

	// Check if the key has expired
	ttl, err := expire.TTL(key)
	assert.Nil(t, err)
	assert.Equal(t, int64(-2), ttl) // -2 indicates that the key does not exist
}

func TestExpireStructure_PEXPIRE(t *testing.T) {
	expire, _ := initExpire()
	defer expire.db.Close()

	key := "test_key"
	milliseconds := int64(2000)

	// Test PEXPIRE function
	expireErr = expire.PEXPIRE(key, milliseconds)
	assert.Nil(t, expireErr)

	// Wait for the key to expire
	time.Sleep(time.Millisecond * time.Duration(milliseconds+100))

	// Check if the key has expired
	ttl, err := expire.TTL(key)
	assert.Nil(t, err)
	assert.Equal(t, int64(-2), ttl) // -2 indicates that the key does not exist
}

func TestExpireStructure_EXPIREAT(t *testing.T) {
	expire, _ := initExpire()
	defer expire.db.Close()

	key := "test_key"
	timestamp := time.Now().Add(time.Second * 2).Unix()

	// Test EXPIREAT function
	expireErr = expire.EXPIREAT(key, timestamp)
	assert.Nil(t, expireErr)

	// Wait for the key to expire
	time.Sleep(time.Second * 3)

	// Check if the key has expired
	ttl, err := expire.TTL(key)
	assert.Nil(t, err)
	assert.Equal(t, int64(-2), ttl) // -2 indicates that the key does not exist
}

func TestExpireStructure_PEXPIREAT(t *testing.T) {
	expire, _ := initExpire()
	defer expire.db.Close()

	key := "test_key"
	timestamp := time.Now().Add(time.Second*2).Unix() * 1000

	// Test PEXPIREAT function
	expireErr = expire.PEXPIREAT(key, timestamp)
	assert.Nil(t, expireErr)

	// Wait for the key to expire
	time.Sleep(time.Second * 3)

	// Check if the key has expired
	ttl, err := expire.TTL(key)
	assert.Nil(t, err)
	assert.Equal(t, int64(-2), ttl) // -2 indicates that the key does not exist
}

func TestExpireStructure_TTL(t *testing.T) {
	expire, _ := initExpire()
	defer expire.db.Close()

	key := "test_key"
	seconds := int64(2)

	// Test TTL function with an existing key
	expire.EXPIRE(key, seconds)
	ttl, err := expire.TTL(key)
	assert.Nil(t, err)
	assert.True(t, ttl > 0) // TTL should be positive

	// Test TTL function with a non-existing key
	ttl, err = expire.TTL("non_existing_key")
	assert.Nil(t, err)
	assert.Equal(t, int64(-2), ttl) // -2 indicates that the key does not exist
}

func TestExpireStructure_PTTL(t *testing.T) {
	expire, _ := initExpire()
	defer expire.db.Close()

	key := "test_key"
	milliseconds := int64(2000)

	// Test PTTL function with an existing key
	expire.PEXPIRE(key, milliseconds)
	ttl, err := expire.PTTL(key)
	assert.Nil(t, err)
	assert.True(t, ttl > 0) // TTL should be positive

	// Test PTTL function with a non-existing key
	ttl, err = expire.PTTL("non_existing_key")
	assert.Nil(t, err)
	assert.Equal(t, int64(-2), ttl) // -2 indicates that the key does not exist
}

func TestExpireStructure_PERSIST(t *testing.T) {
	expire, _ := initExpire()
	defer expire.db.Close()

	key := "test_key"
	seconds := int64(2)

	// Set an expiration time for the key
	expire.EXPIRE(key, seconds)

	// Test PERSIST function
	expireErr = expire.PERSIST(key, seconds)
	assert.Nil(t, expireErr)

	// Check if the key is still in the database and has no expiration time
	ttl, err := expire.TTL(key)
	assert.Nil(t, err)
	assert.Equal(t, int64(-1), ttl) // -1 indicates that the key has no expiration
}

func TestExpireStructure_EXPIREBY(t *testing.T) {
	expire, _ := initExpire()
	defer expire.db.Close()

	key := "test_key"
	duration := int64(2)

	// Test EXPIREBY function without tag
	expireErr = expire.EXPIREBY(key, duration)
	assert.Nil(t, expireErr)

	// Wait for the key to expire
	time.Sleep(time.Second * time.Duration(duration+1))

	// Check if the key has expired
	ttl, err := expire.TTL(key)
	assert.Nil(t, err)
	assert.Equal(t, int64(-2), ttl) // -2 indicates that the key does not exist

	// Test EXPIREBY function with tag "ex"
	expireErr = expire.EXPIREBY(key, duration, "ex")
	assert.Nil(t, expireErr)

	// Wait for the key to expire
	time.Sleep(time.Second * time.Duration(duration+1))

	// Check if the key has expired
	ttl, err = expire.TTL(key)
	assert.Nil(t, err)
	assert.Equal(t, int64(-2), ttl) // -2 indicates that the key does not exist

	// Test EXPIREBY function with an invalid tag
	expireErr = expire.EXPIREBY(key, duration, "invalid_tag")
	assert.Equal(t, ErrInvalidArgs, expireErr)
}

func TestExpireStructure_PEXPIREBY(t *testing.T) {
	expire, _ := initExpire()
	defer expire.db.Close()

	key := "test_key"
	duration := int64(2000)

	// Test PEXPIREBY function without tag
	expireErr = expire.PEXPIREBY(key, duration)
	assert.Nil(t, expireErr)

	// Wait for the key to expire
	time.Sleep(time.Millisecond * time.Duration(duration+100))

	// Check if the key has expired
	ttl, err := expire.TTL(key)
	assert.Nil(t, err)
	assert.Equal(t, int64(-2), ttl) // -2 indicates that the key does not exist

	// Test PEXPIREBY function with tag "ex"
	expireErr = expire.PEXPIREBY(key, duration, "ex")
	assert.Nil(t, expireErr)

	// Wait for the key to expire
	time.Sleep(time.Millisecond * time.Duration(duration+100))

	// Check if the key has expired
	ttl, err = expire.TTL(key)
	assert.Nil(t, err)
	assert.Equal(t, int64(-2), ttl) // -2 indicates that the key does not exist

	// Test PEXPIREBY function with an invalid tag
	expireErr = expire.PEXPIREBY(key, duration, "invalid_tag")
	assert.Equal(t, ErrInvalidArgs, expireErr)
}

func TestExpireStructure_EXPIREBYAT(t *testing.T) {
	expire, _ := initExpire()
	defer expire.db.Close()

	key := "test_key"
	timestamp := time.Now().Add(time.Second * 2).Unix()

	// Test EXPIREBYAT function without tag
	expireErr = expire.EXPIREBYAT(key, timestamp)
	assert.Nil(t, expireErr)

	// Wait for the key to expire
	time.Sleep(time.Second * 3)

	// Check if the key has expired
	ttl, err := expire.TTL(key)
	assert.Nil(t, err)
	assert.Equal(t, int64(-2), ttl) // -2 indicates that the key does not exist

	// Test EXPIREBYAT function with tag "ex"
	expireErr = expire.EXPIREBYAT(key, timestamp, "ex")
	assert.Nil(t, expireErr)

	// Check if the key has expired
	ttl, err = expire.TTL(key)
	assert.Nil(t, err)
	assert.Equal(t, int64(timestamp), ttl)

	// Test EXPIREBYAT function with an invalid tag
	expireErr = expire.EXPIREBYAT(key, timestamp, "invalid_tag")
	assert.Equal(t, ErrInvalidArgs, expireErr)
}

func TestExpireStructure_PEXPIREBYAT(t *testing.T) {
	expire, _ := initExpire()
	defer expire.db.Close()

	key := "test_key"
	timestamp := time.Now().Add(time.Second*2).Unix() * 1000

	// Test PEXPIREBYAT function without tag
	expireErr = expire.PEXPIREBYAT(key, timestamp)
	assert.Nil(t, expireErr)

	// Wait for the key to expire
	time.Sleep(time.Second * 3)

	// Check if the key has expired
	ttl, err := expire.TTL(key)
	assert.Nil(t, err)
	assert.Equal(t, int64(-2), ttl) // -2 indicates that the key does not exist

	// Test PEXPIREBYAT function with tag "ex"
	expireErr = expire.PEXPIREBYAT(key, timestamp, "ex")
	assert.Nil(t, expireErr)

	// Check if the key has expired
	ttl, err = expire.TTL(key)
	assert.Nil(t, err)
	assert.Equal(t, int64(timestamp/1000), ttl) // -2 indicates that the key does not exist

	// Test PEXPIREBYAT function with an invalid tag
	expireErr = expire.PEXPIREBYAT(key, timestamp, "invalid_tag")
	assert.Equal(t, ErrInvalidArgs, expireErr)
}
