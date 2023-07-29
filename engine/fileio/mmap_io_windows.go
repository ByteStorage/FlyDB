//go:build windows

package fileio

import (
	"errors"
	"os"
	"sync/atomic"
	"syscall"
	"unsafe"
)

type MMapIO struct {
	fd       *os.File // system file descriptor
	handle   syscall.Handle
	data     []byte // the mapping area corresponding to the file
	offset   int64  // next write location
	fileSize int64  // max file size
	fileName string
	count    atomic.Uint32 // the count of dbs using this mmap io
}

func (mio *MMapIO) Init() (*MMapIO, error) {
	fd, err := os.OpenFile(mio.fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, DataFilePerm)
	if err != nil {
		return nil, err
	}
	mio.fd = fd

	info, _ := fd.Stat()
	mio.offset = info.Size()

	// Expand files to maximum file size, crop when saving
	if err := fd.Truncate(mio.fileSize); err != nil {
		return nil, err
	}

	// Building mappings between memory and disk files
	h, err := syscall.CreateFileMapping(syscall.Handle(mio.fd.Fd()), nil,
		syscall.PAGE_READWRITE, 0, uint32(mio.fileSize), nil)
	if err != nil {
		return nil, err
	}
	mio.handle = h

	addr, _ := syscall.MapViewOfFile(h, syscall.FILE_MAP_WRITE, 0,
		0, uintptr(mio.fileSize))
	if err != nil {
		return nil, err
	}
	mio.data = *(*[]byte)(unsafe.Pointer(addr))

	return mio, nil
}

func (mio *MMapIO) unmap() error {
	if mio.data == nil {
		return nil
	}

	addr := (uintptr)(unsafe.Pointer(&mio.data[0]))
	err := syscall.UnmapViewOfFile(addr)
	mio.data = nil

	return err
}

// Read Copy data from the mapping area to byte slice
func (mio *MMapIO) Read(b []byte, offset int64) (int, error) {
	return copy(b, mio.data[offset:]), nil
}

// Write Copy data from byte slice to the mapping area
func (mio *MMapIO) Write(b []byte) (int, error) {
	oldOffset := mio.offset
	newOffset := mio.offset + int64(len(b))
	if newOffset > mio.fileSize {
		return 0, errors.New("exceed file max content length")
	}

	mio.offset = newOffset
	return copy(mio.data[oldOffset:newOffset], b), nil
}

func (mio *MMapIO) Sync() error {
	err := syscall.FlushFileBuffers(mio.handle)
	if err != nil {
		return err
	}
	return nil
}

// Size return the size of current file
func (mio *MMapIO) Size() (int64, error) {
	return mio.offset, nil
}

func (mio *MMapIO) Close() (err error) {
	controller.lock.Lock()
	defer controller.lock.Unlock()

	mio.SubCount()
	if mio.GetCount() > 0 {
		return nil
	}

	delete(controller.files, mio.fileName)

	if err = mio.fd.Truncate(mio.offset); err != nil {
		return err
	}
	if err = mio.Sync(); err != nil {
		return err
	}
	if err = mio.unmap(); err != nil {
		panic(err)
	}
	return syscall.CloseHandle(mio.handle)
}

func (mio *MMapIO) GetCount() uint32 {
	return mio.count.Load()
}

func (mio *MMapIO) AddCount() {
	mio.count.Add(1)
}

func (mio *MMapIO) SubCount() {
	mio.count.Add(-1)
}
