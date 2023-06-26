package wal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	wal, err := New()
	assert.Nil(t, err)
	assert.NotNil(t, wal)
	index, err := wal.log.LastIndex()
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), index)
}

func TestWal_Write(t *testing.T) {
	wal, err := New()
	assert.Nil(t, err)
	index, err := wal.log.LastIndex()

	assert.Nil(t, err)
	assert.Equal(t, uint64(1), index)

	data := []byte("test data")
	err = wal.Write(data)
	assert.Nil(t, err)

	index, err = wal.log.LastIndex()
	assert.Nil(t, err)
	assert.Equal(t, uint64(2), index)
}

func TestWal_Read(t *testing.T) {
	wal, err := New()
	assert.Nil(t, err)

	data := []byte("test data")
	err = wal.Write(data)
	assert.Nil(t, err)

	index, err := wal.log.LastIndex()

	readData, err := wal.Read(index)
	assert.Nil(t, err)
	assert.Equal(t, data, readData)
}

func TestWal_ReadLast(t *testing.T) {
	wal, err := New()
	assert.Nil(t, err)

	data := []byte("test data")
	err = wal.Write(data)
	assert.Nil(t, err)

	readData, err := wal.ReadLast()
	assert.Nil(t, err)
	assert.Equal(t, data, readData)
}
