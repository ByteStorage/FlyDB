//go:build linux

package fileio

import (
	"errors"
	atomic2 "go.uber.org/atomic"
	"os"
	"syscall"
	"unsafe"
)

type MMapIO struct {
	fd       *os.File // system file descriptor
	data     []byte   // the mapping area corresponding to the file
	offset   int64    // next write location
	fileSize int64    // max file size
	fileName string
	count    atomic2.Int32 // the count of dbs using this mmap io
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
	b, err := syscall.Mmap(int(mio.fd.Fd()), 0, int(mio.fileSize),
		syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	mio.data = b

	return mio, nil
}

// UnMapUnix Unmapping between memory and files
func (mio *MMapIO) unmap() error {
	if mio.data == nil {
		return nil
	}

	err := syscall.Munmap(mio.data)
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

// Sync Synchronize data from memory to disk
func (mio *MMapIO) Sync() error {
	_, _, err := syscall.Syscall(syscall.SYS_MSYNC, uintptr(unsafe.Pointer(&mio.data[0])), uintptr(mio.offset), uintptr(syscall.MS_SYNC))
	if err != 0 {
		return err
	}
	return nil
}

// Size return the size of current file
func (mio *MMapIO) Size() (int64, error) {
	return mio.offset, nil
}

// Close file
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
	return mio.fd.Close()
}

func (mio *MMapIO) GetCount() int32 {
	return mio.count.Load()
}

func (mio *MMapIO) AddCount() {
	mio.count.Add(1)
}

func (mio *MMapIO) SubCount() {
	mio.count.Add(-1)
}
