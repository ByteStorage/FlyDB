package service

import (
	"github.com/ByteStorage/FlyDB/lib/proto/ghash"
	"github.com/ByteStorage/FlyDB/lib/proto/gstring"
	"github.com/ByteStorage/FlyDB/structure"
	"os"
)

// Service is a grpc Service for db
type Service struct {
	gstring.GStringServiceServer
	ghash.GHashServiceServer
	Addr string // db server address
	dbs  *structure.StringStructure
	dbh  *structure.HashStructure
	sig  chan os.Signal
}

// NewService returns a new grpc Service
func NewService(addr string, db *structure.StringStructure) *Service {
	return &Service{
		Addr: addr,
		dbs:  db,
		sig:  make(chan os.Signal),
	}
}

// NewHashService returns a new grpc Service
func NewHashService(addr string, db *structure.HashStructure) *Service {
	return &Service{
		Addr: addr,
		dbh:  db,
		sig:  make(chan os.Signal),
	}
}
