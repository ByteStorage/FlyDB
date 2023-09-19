package bloom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBloomFilter(t *testing.T) {
	filter := NewBloomFilter(1000, 0.01)
	assert.NotNil(t, filter)
}

func TestFilter_Add(t *testing.T) {
	filter := NewBloomFilter(1000, 0.01)
	filter.Add([]byte("hello"))
	assert.True(t, filter.MayContainItem([]byte("hello")))
	assert.False(t, filter.MayContainItem([]byte("world")))
}

func TestFilter_MayContainItem(t *testing.T) {
	filter := NewBloomFilter(1000, 0.01)
	filter.Add([]byte("hello"))
	filter.Add([]byte("world"))
	filter.Add([]byte("flydb"))
	filter.Add([]byte("bloom"))
	assert.True(t, filter.MayContainItem([]byte("hello")))
	assert.True(t, filter.MayContainItem([]byte("world")))
	assert.True(t, filter.MayContainItem([]byte("flydb")))
	assert.True(t, filter.MayContainItem([]byte("bloom")))
	assert.False(t, filter.MayContainItem([]byte("fly")))
	assert.False(t, filter.MayContainItem([]byte("db")))
}
