package fileio

import (
	"sync"
)

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
		v.AddCount()
		return v, nil
	}

	manager, err := (&MMapIO{fileName: fileName, fileSize: fileSize}).Init()
	return manager, err
}
