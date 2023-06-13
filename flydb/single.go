package flydb

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
)

// NewFlyDB open a new db instance
func NewFlyDB(options config.Options) (*engine.DB, error) {
	return engine.NewDB(options)
}
