package memory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMemTable(t *testing.T) {
	table := NewMemTable()
	assert.NotNil(t, table)
}

func TestMemTable_Get(t *testing.T) {
	table := NewMemTable()
	assert.NotNil(t, table)

	value, err := table.Get("test")
	assert.Nil(t, value)
	assert.NotNil(t, err)
	assert.Equal(t, "key not found", err.Error())
}

func TestMemTable_Put(t *testing.T) {
	table := NewMemTable()
	assert.NotNil(t, table)

	table.Put("test", []byte("test"))
	value, err := table.Get("test")
	assert.Nil(t, err)
	assert.NotNil(t, value)
	assert.Equal(t, "test", string(value))
}
