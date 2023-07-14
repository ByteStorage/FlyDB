package service

import (
	"github.com/ByteStorage/FlyDB/lib/proto/gstring"
	"github.com/ByteStorage/FlyDB/structure"
	"os"
)

// Service is a grpc Service for db
type Service struct {
	gstring.GStringServiceServer
	Addr string // db server address
	db   *structure.StringStructure
	sig  chan os.Signal
}

// NewService returns a new grpc Service
func NewService(addr string, db *structure.StringStructure) *Service {
	return &Service{
		Addr: addr,
		db:   db,
		sig:  make(chan os.Signal),
	}
}
