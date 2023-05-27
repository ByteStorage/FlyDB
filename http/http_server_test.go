package http_handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qishenonly/flydb"
)

func newHttpHandler() (*HttpHandler, error) {
	options := flydb.DefaultOptions
	options.DirPath = "/tmp/flydb"
	db, err := flydb.NewFlyDB(options)
	if err != nil {
		return nil, err
	}
	Serve := NewHttpHandler(db)
	return Serve, nil
}

func TestNewHTTPServer(t *testing.T) {
	server, err := newHttpHandler()
	defer func(server *HttpHandler) {
		err := server.Close()
		if err != nil {

		}
	}(server)
	if err != nil {
		t.Error(err)
	}
}

// 测试Put方法
func TestPut(t *testing.T) {
	handler, _ := newHttpHandler()
	// 创建一个测试用的http server
	server := httptest.NewServer(http.HandlerFunc(handler.PutHandler))
	defer server.Close()
	// 构造请求
	reqBody := map[string]string{
		"key":   "test_key",
		"value": "test_value",
	}
	reqBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, server.URL, bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 检查响应
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ReadAll error: %v", err)
	}
	if string(body) != "ok" {
		t.Errorf("Put error: expected ok, got %s", string(body))
	}

	// 验证是否put成功
	valueByte, err := handler.Get([]byte("test_key"))
	if err != nil {
		t.Errorf("Get error: %v", err)
	}
	value := string(valueByte)
	if value != "test_value" {
		t.Errorf("Put error: expected %s, got %s", "test_value", value)
	} else {
		t.Logf("Put: test_value and Get: %s", value)
	}
}
