package fio

import (
	"golang.org/x/sys/unix"
	"os"
	"sync/atomic"
	"syscall"
)

const size = 1024 * 1024 * 256

type MmapIO struct {
	fd     *os.File //系统文件描述符
	data   []byte
	offset int64 // 用于记录写入位置
}

func NewMmapIOManager(filePath string) (*MmapIO, error) {
	fd, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, DataFilePerm)
	if err := fd.Truncate(size); err != nil {
		return nil, err
	}
	data, err := syscall.Mmap(int(fd.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	return &MmapIO{
		fd:     fd,
		data:   data,
		offset: 0,
	}, nil
}

func (fd *MmapIO) Read(data []byte, pos int64) (int, error) {
	start := int(pos)
	end := start + len(data)
	if start >= len(fd.data) {
		return 0, nil
	}
	if end > len(fd.data) {
		end = len(fd.data)
	}
	return copy(data, fd.data[start:end]), nil
}

func (fd *MmapIO) Write(data []byte) (int, error) {
	currentOffset := atomic.LoadInt64(&fd.offset)        // 获取当前写入位置
	requiredCapacity := currentOffset + int64(len(data)) // 计算所需的容量

	// 如果 fd.data 的容量不足以容纳所需的容量，进行扩容
	if int64(cap(fd.data)) < requiredCapacity {
		newCapacity := requiredCapacity + 1024 // 增加一些额外的容量
		newData := make([]byte, requiredCapacity, newCapacity)
		copy(newData, fd.data)
		fd.data = newData
	}

	newData := append(fd.data[:currentOffset], data...)    // 将新数据追加到原有数据之后
	copy(fd.data[currentOffset:], newData[currentOffset:]) // 将追加后的数据复制到 fd.data
	atomic.AddInt64(&fd.offset, int64(len(data)))          // 更新写入位置
	return len(data), nil
}

func (fd *MmapIO) Sync() error {
	return unix.Msync(fd.data, syscall.MS_SYNC)
}

func (fd *MmapIO) Close() error {
	if err := syscall.Munmap(fd.data); err != nil {
		return err
	}
	return fd.fd.Close()
}

func (fd *MmapIO) Size() (int64, error) {
	fileInfo, err := fd.fd.Stat()
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func NewIOManager(filename string) (IOManager, error) {
	return NewMmapIOManager(filename)
}
