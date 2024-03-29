package fileio

import (
	"errors"
	"github.com/edsrzf/mmap-go"
	"os"
)

type MMapIO struct {
	fd       *os.File  // system file descriptor
	data     mmap.MMap // the mapping area corresponding to the file
	dirty    bool      // has changed
	offset   int64     // next write location
	fileSize int64     // max file size
}

// NewMMapIOManager Initialize Mmap IO
func NewMMapIOManager(fileName string, fileSize int64) (*MMapIO, error) {
	mmapIO := &MMapIO{fileSize: fileSize}

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
	b, err := mmap.Map(fd, mmap.RDWR, 0)
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
	mio.dirty = true
	return copy(mio.data[oldOffset:], b), nil
}

// Sync Synchronize data from memory to disk
func (mio *MMapIO) Sync() error {
	if !mio.dirty {
		return nil
	}

	if err := mio.data.Flush(); err != nil {
		return err
	}

	mio.dirty = false
	return nil
}

// Close file
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

// Size return the size of current file
func (mio *MMapIO) Size() (int64, error) {
	return mio.offset, nil
}

// UnMap Unmapping between memory and files
func (mio *MMapIO) UnMap() error {
	if mio.data == nil {
		return nil
	}
	err := mio.data.Unmap()
	mio.data = nil
	return err
}
