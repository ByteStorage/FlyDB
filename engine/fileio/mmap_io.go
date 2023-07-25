package fileio

import (
	"errors"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

type MMapIO struct {
	fd       *os.File // system file descriptor
	data     []byte   // the mapping area corresponding to the file
	offset   int64    // next write location
	fileSize int64    // max file size
	fileName string
	count    int // the count of dbs using this mmap io
}

type mmapFileController struct {
	lock  sync.Mutex
	files map[string]*MMapIO
}

var controller = mmapFileController{
	lock:  sync.Mutex{},
	files: map[string]*MMapIO{},
}

// NewMMapIOManager Initialize Mmap IO
func NewMMapIOManager(fileName string, fileSize int64) (*MMapIO, error) {
	controller.lock.Lock()
	defer controller.lock.Unlock()

	if v, ok := controller.files[fileName]; ok {
		v.count++
		return v, nil
	}

	mmapIO := &MMapIO{fileSize: fileSize, fileName: fileName, count: 1}
	controller.files[fileName] = mmapIO

	fd, err := os.OpenFile(
		fileName,
		os.O_CREATE|os.O_RDWR|os.O_APPEND,
		DataFilePerm,
	)
	if err != nil {
		return nil, err
	}
	info, _ := fd.Stat()

	// Expand files to maximum file size, crop when saving
	if err := fd.Truncate(fileSize); err != nil {
		return nil, err
	}

	// Building mappings between memory and disk files
	b, err := syscall.Mmap(int(fd.Fd()), 0, int(fileSize), syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	mmapIO.fd = fd
	mmapIO.data = b
	mmapIO.offset = info.Size()
	return mmapIO, nil
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

// Close file
func (mio *MMapIO) Close() (err error) {
	controller.lock.Lock()
	defer controller.lock.Unlock()

	mio.count--
	if mio.count > 0 {
		return nil
	}

	delete(controller.files, mio.fileName)

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

// Size return the size of current file
func (mio *MMapIO) Size() (int64, error) {
	return mio.offset, nil
}

// UnMap Unmapping between memory and files
func (mio *MMapIO) UnMap() error {
	if mio.data == nil {
		return nil
	}
	err := syscall.Munmap(mio.data)
	mio.data = nil
	return err
}
