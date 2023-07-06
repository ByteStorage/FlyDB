package client

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	"github.com/ByteStorage/FlyDB/engine/grpc/service"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestClient_Put(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-put")
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := engine.NewDB(opts)
	defer db.Clean()
	assert.Nil(t, err)
	s := service.NewService(config.DefaultAddr, db)
	go s.StartServer()
	//wait for server start
	for {
		if s.IsGrpcServerRunning() {
			break
		}
	}
	client := &Client{
		Addr: config.DefaultAddr,
	}
	err = client.Put([]byte("test"), []byte("test"))
	assert.Nil(t, err)
}
