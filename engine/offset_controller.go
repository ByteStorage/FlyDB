package engine

import "sync"

type OffsetController struct {
	m      sync.Mutex
	offset map[uint32]int64
}

var Controller = &OffsetController{
	m:      sync.Mutex{},
	offset: map[uint32]int64{},
}

func SingleOffset() *OffsetController {
	return Controller
}

func (c *OffsetController) CanWrite(fileId uint32, filesize int64, size int64) (int64, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	if v, ok := c.offset[fileId]; ok {
		if v+size <= filesize {
			c.offset[fileId] = v + size
			return v, true
		}
	}
	return 0, false
}

func (c *OffsetController) AddNew(fileId uint32, offset int64) int64 {
	c.m.Lock()
	defer c.m.Unlock()
	if _, ok := c.offset[fileId]; ok {
		return c.offset[fileId]
	}
	c.offset[fileId] = offset
	return offset
}

func (c *OffsetController) ChangeOffset(fileId uint32, offset int64) {
	c.m.Lock()
	defer c.m.Unlock()
	c.offset[fileId] = offset
}
