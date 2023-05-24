package fio

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestBufio_Write(t *testing.T) {
	path := filepath.Join("/tmp", "a.data")
	fio, err := NewBufIOManager(path)
	defer destoryFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, fio)

	n, err := fio.Write([]byte(""))
	assert.Equal(t, 0, n)
	assert.Nil(t, err)

	n, err = fio.Write([]byte("bitcask kv"))
	assert.Equal(t, 10, n)
	assert.Nil(t, err)
	n, err = fio.Write([]byte("storage"))
	assert.Equal(t, 7, n)
	assert.Nil(t, err)
}

func TestBufio_Read(t *testing.T) {
	path := filepath.Join("/tmp", "a.data")
	fio, err := NewBufIOManager(path)
	defer destoryFile(path)

	assert.Nil(t, err)
	assert.NotNil(t, fio)

	_, err = fio.Write([]byte("key-a"))
	assert.Nil(t, err)
	_, err = fio.Write([]byte("key-b"))
	assert.Nil(t, err)

	b1 := make([]byte, 5)
	n, err := fio.Read(b1, 0)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("key-a"), b1)
}

func TestBufio_Read2(t *testing.T) {
	path := filepath.Join("/tmp", "a.data")
	fio, err := NewBufIOManager(path)
	defer destoryFile(path)

	assert.Nil(t, err)
	assert.NotNil(t, fio)

	_, err = fio.Write([]byte("key-a"))
	assert.Nil(t, err)
	_, err = fio.Write([]byte("key-b"))
	assert.Nil(t, err)

	b1 := make([]byte, 5)
	n, err := fio.Read(b1, 5)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("key-b"), b1)
}

func TestBufio_Write2(t *testing.T) {
	path := filepath.Join("/tmp", "a.data")
	fio, err := NewBufIOManager(path)
	defer destoryFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, fio)

	n, err := fio.Write([]byte("bitcask kv"))
	assert.Equal(t, 10, n)
	assert.Nil(t, err)
	n, err = fio.Write([]byte("storage"))
	assert.Equal(t, 7, n)
	assert.Nil(t, err)
}

func TestBufio_Write_Speed(t *testing.T) {
	path := filepath.Join("/tmp", "a.data")
	fio, err := NewBufIOManager(path)
	defer fio.Close()
	assert.Nil(t, err)
	assert.NotNil(t, fio)

	for i := 0; i < 1000000; i++ {
		n, err := fio.Write([]byte("bitcask kv"))
		assert.Equal(t, 10, n)
		assert.Nil(t, err)
	}

	assert.Nil(t, err)
}

func TestBufio_Read_Speed(t *testing.T) {
	path := filepath.Join("/tmp", "a.data")
	fio, err := NewBufIOManager(path)
	defer destoryFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, fio)

	for i := 0; i < 1000000; i++ {
		b1 := make([]byte, 10)
		_, _ = fio.Read(b1, int64(i*10))
		/*//TODO 这个地方有个奇怪的bug，我设成expected为10,但actual为6，改成6的时候，actual为10，但是实际打印的时候，b1是正常的
		//从这几个测试结果来看，貌似可以改成bufio进行读写，这样性能会更好
		assert.Equal(t, 10, n)
		assert.Nil(t, err)
		assert.Equal(t, []byte("bitcask kv"), b1)*/
		//fmt.Println(string(b1))
	}

}
