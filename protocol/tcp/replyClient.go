package tcp

import (
	"github.com/ByteStorage/FlyDB/lib/sync/wait"
	"net"
	"time"
)

//var _ tcpIF.Handler = (*ReplyClient)(nil)

type ReplyClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

// Close client connection
func (r *ReplyClient) Close() error {
	r.Waiting.WaitTimeout(10 * time.Second)
	_ = r.Conn.Close()
	return nil
}
