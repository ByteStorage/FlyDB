package tcp

import (
	"bufio"
	"context"
	"fmt"
	"github.com/ByteStorage/flydb/protocol/tcpIF"
	"github.com/ByteStorage/flydb/sync/boolAm"
	"io"
	"net"
	"sync"
)

var _ tcpIF.Handler = (*TcpReplyHandler)(nil)

type TcpReplyHandler struct {
	activeConn sync.Map
	isClosed   boolAm.Boolean
}

// NewHandler create a new handler
func NewHandler() *TcpReplyHandler {
	return &TcpReplyHandler{}
}

// Handle client connection
func (t *TcpReplyHandler) Handle(ctx context.Context, conn net.Conn) {
	if t.isClosed.GetBoolAtomic() {
		_ = conn.Close()
		return
	}

	client := &ReplyClient{
		Conn: conn,
	}
	t.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("client close connection")
				t.activeConn.Delete(client)
			} else {
				fmt.Println("read message error: ", err)
			}
			return
		}

		if msg[:2] == "\\n" {
			client.Waiting.Add(1)
			buf := []byte(msg)
			_, _ = conn.Write(buf)
			client.Waiting.Done()
		}
	}

}

// Close handler
func (t *TcpReplyHandler) Close() error {
	fmt.Println("tcp server close")
	t.isClosed.SetBoolAtomic(true)
	t.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*ReplyClient)
		_ = client.Close()
		return true
	})
	return nil
}
