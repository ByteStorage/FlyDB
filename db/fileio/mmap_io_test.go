package fileio

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestNewMMapIOManager(t *testing.T) {
	path := filepath.Join("/tmp", "a.data")
	mio, err := NewMMapIOManager(path, DefaultFileSize)
	defer destoryFile(path)

	assert.Nil(t, err)
	assert.NotNil(t, mio)
}

func TestMMapIO_Write(t *testing.T) {
	path := filepath.Join("/tmp", "a.data")
	mio, err := NewMMapIOManager(path, DefaultFileSize)
	defer destoryFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, mio)

	n, err := mio.Write([]byte(""))
	assert.Equal(t, 0, n)
	assert.Nil(t, err)

	n, err = mio.Write([]byte("bitcask kv"))
	assert.Equal(t, 10, n)
	assert.Nil(t, err)
	n, err = mio.Write([]byte("storage"))
	assert.Equal(t, 7, n)
	assert.Nil(t, err)
}

func TestMMapIO_Read(t *testing.T) {
	path := filepath.Join("/tmp", "a.data")
	mio, err := NewMMapIOManager(path, DefaultFileSize)
	defer destoryFile(path)

	assert.Nil(t, err)
	assert.NotNil(t, mio)

	_, err = mio.Write([]byte("key-a"))
	assert.Nil(t, err)
	_, err = mio.Write([]byte("key-b"))
	assert.Nil(t, err)

	b1 := make([]byte, 5)
	n, err := mio.Read(b1, 0)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("key-a"), b1)

	b2 := make([]byte, 5)
	n, err = mio.Read(b2, 5)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("key-b"), b2)
}

func TestMMapIO_Sync(t *testing.T) {
	path := filepath.Join("/tmp", "a.data")
	mio, err := NewMMapIOManager(path, DefaultFileSize)
	defer destoryFile(path)

	assert.Nil(t, err)
	assert.NotNil(t, mio)

	_, err = mio.Write([]byte("key-c"))
	err = mio.Sync()
	assert.Nil(t, err)
}

func TestMMapIO_Close(t *testing.T) {
	path := filepath.Join("/tmp", "a.data")
	mio, err := NewMMapIOManager(path, DefaultFileSize)
	defer destoryFile(path)

	assert.Nil(t, err)
	assert.NotNil(t, mio)

	_, err = mio.Write([]byte("key-a"))
	assert.Nil(t, err)

	err = mio.Close()
	assert.Nil(t, err)
}

func TestMMapIO_Write_Speed(t *testing.T) {
	path := filepath.Join("/tmp", "a.data")
	mio, err := NewMMapIOManager(path, DefaultFileSize)
	assert.Nil(t, err)
	assert.NotNil(t, mio)

	for i := 0; i < 1000000; i++ {
		n, err := mio.Write([]byte("bitcask kv"))
		assert.Equal(t, 10, n)
		assert.Nil(t, err)
	}

	err = mio.Close()
	assert.Nil(t, err)
}

func TestMMapIO_Read_Speed(t *testing.T) {
	path := filepath.Join("/tmp", "a.data")
	mio, err := NewMMapIOManager(path, DefaultFileSize)
	defer destoryFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, mio)

	for i := 0; i < 1000000; i++ {
		b1 := make([]byte, 10)
		n, err := mio.Read(b1, int64(i*10))
		assert.Equal(t, 10, n)
		assert.Nil(t, err)
		//assert.Equal(t, []byte("bitcask kv"), b1)
	}
}
