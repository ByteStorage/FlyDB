package flydb

import (
	"github.com/ByteStorage/flydb/config"
	"github.com/ByteStorage/flydb/engine"
)

// NewFlyDB open a new db instance
func NewFlyDB(options config.Options) (*engine.DB, error) {
	return engine.NewDB(options)
}
