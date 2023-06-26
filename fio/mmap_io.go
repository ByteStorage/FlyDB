package fio

import (
	"errors"
	"os"
	"syscall"
	"unsafe"
)

// DefaultMemMapSize 最大映射内存大小
const DefaultMemMapSize = 256 * 1024 * 1024

// MMapIO 标准系统文件IO
type MMapIO struct {
	fd     *os.File // 系统文件描述符
	data   []byte   // 与文件对应的映射区
	dirty  bool     // 是否更改过
	offset int64    // 写入位置
}

// NewMMapIOManager 初始化标准文件 IO
func NewMMapIOManager(fileName string) (*MMapIO, error) {
	mmapIO := &MMapIO{}
	fd, err := os.OpenFile(
		fileName,
		os.O_CREATE|os.O_RDWR|os.O_APPEND,
		DataFilePerm,
	)
	if err != nil {
		return nil, err
	}
	info, _ := fd.Stat()

	// 将文件扩容到映射区大小, 保存时会裁剪
	if err := fd.Truncate(DefaultMemMapSize); err != nil {
		return nil, err
	}

	// 构建映射
	b, err := syscall.Mmap(int(fd.Fd()), 0, DefaultMemMapSize, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	mmapIO.fd = fd
	mmapIO.data = b
	mmapIO.offset = info.Size()
	return mmapIO, nil
}

func (mio *MMapIO) Read(b []byte, offset int64) (int, error) {
	return copy(b, mio.data[offset:]), nil
}

func (mio *MMapIO) Write(b []byte) (int, error) {
	oldOffset := mio.offset
	newOffset := mio.offset + int64(len(b))
	if newOffset > DefaultMemMapSize {
		return 0, errors.New("exceed file max content length")
	}

	mio.offset = newOffset
	mio.dirty = true
	return copy(mio.data[oldOffset:], b), nil
}

func (mio *MMapIO) Sync() error {
	if !mio.dirty {
		return nil
	}

	_, _, err := syscall.Syscall(syscall.SYS_MSYNC, uintptr(unsafe.Pointer(&mio.data[0])), uintptr(mio.offset), uintptr(syscall.MS_SYNC))
	if err != 0 {
		return err
	}

	mio.dirty = false
	return nil
}

func (mio *MMapIO) Close() (err error) {
	if err = mio.fd.Truncate(mio.offset); err != nil {
		return err
	}
	if err = mio.Sync(); err != nil {
		return err
	}
	if err = mio.UnMap(); err != nil {
		panic(err)
	}
	return mio.fd.Close()
}

func (mio *MMapIO) Size() (int64, error) {
	return mio.offset, nil
}

func (mio *MMapIO) UnMap() error {
	if mio.data == nil {
		return nil
	}
	err := syscall.Munmap(mio.data)
	mio.data = nil
	return err
}
