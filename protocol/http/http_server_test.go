package http_handler

import (
	"bytes"
	"encoding/json"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newHttpHandler() (*HttpHandler, error) {
	options := config.DefaultOptions
	options.DirPath = "/tmp/flydb"
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}
	Serve := NewHttpHandler(db)
	return Serve, nil
}

func TestNewHTTPServer(t *testing.T) {
	server, err := newHttpHandler()
	time.Sleep(time.Millisecond * 100)
	defer server.Clean()
	assert.Nil(t, err)
}

// 测试Put方法
func TestPut(t *testing.T) {
	handler, _ := newHttpHandler()
	defer handler.Clean()
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
	if value == "test_value" {
		t.Logf("Put: test_value and Get: %s", value)
	} else {
		t.Errorf("Put error: expected %s, got %s", "test_value", value)
	}
}

func TestDel(t *testing.T) {
	handler, _ := newHttpHandler()
	defer handler.Clean()
	// 创建一个测试用的http server
	server := httptest.NewServer(http.HandlerFunc(handler.DelHandler))
	defer server.Close()
	req, _ := http.NewRequest(http.MethodDelete, server.URL+"?key=test_key", nil)
	req.Header.Set("Content-Type", "multipart/form-data")

	//提前插入test_key
	err := handler.Put([]byte("test_key"), []byte("test_value"))
	if err != nil {
		return
	}

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
		t.Errorf("Delete error: expected ok, got %s", string(body))
	}

	// 验证是否delete成功
	_, err = handler.Get([]byte("test_key"))
	if err == nil {
		t.Errorf("Del error: %v", err)
	} else {
		t.Logf("delete data success")
	}
}

func TestGet(t *testing.T) {
	handler, _ := newHttpHandler()
	defer handler.Clean()
	// 创建一个测试用的http server
	server := httptest.NewServer(http.HandlerFunc(handler.GetHandler))
	defer server.Close()
	req, _ := http.NewRequest(http.MethodGet, server.URL+"?key=test_key", nil)
	req.Header.Set("Content-Type", "application/json")
	//提前插入test_key
	err := handler.Put([]byte("test_key"), []byte("test_value"))
	if err != nil {
		return
	}

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
	// 验证是否Get成功
	val, err := handler.Get([]byte("test_key"))
	if err != nil {
		t.Errorf("Get error: %v", err)
	}

	if string(body) != string(val) {
		t.Errorf("Get error: expected ok, got %s", string(body))
	} else {
		t.Logf("value:%s", string(val))
	}
}

func TestPost(t *testing.T) {
	handler, _ := newHttpHandler()
	defer handler.Clean()
	// 创建一个测试用的 HTTP 服务器
	server := httptest.NewServer(http.HandlerFunc(handler.PostHandler))
	defer server.Close()

	// 构造请求
	reqBody := map[string]string{
		"key":   "test_post",
		"value": "test_post_value",
	}
	reqBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// 处理关闭响应主体失败的错误
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
		t.Errorf("Post error: expected ok, got %s", string(body))
	}

	// 验证是否成功创建新资源
	valueByte, err := handler.Get([]byte("test_post"))
	if err != nil {
		t.Errorf("Get error: %v", err)
	}
	value := string(valueByte)
	if value != "test_post_value" {
		t.Errorf("Post error: expected %s, got %s", "test_value", value)
	} else {
		t.Logf("Post: test_value and Get: %s", value)
	}
}

func TestGetListKeysHandler(t *testing.T) {
	handler, _ := newHttpHandler()
	defer handler.DB.Clean()
	// 创建一个测试用的http server
	server := httptest.NewServer(http.HandlerFunc(handler.GetListKeysHandler))
	defer server.Close()
	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	// 发送请求
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
	// 解析JSON响应
	var keys [][]byte
	err = json.Unmarshal(body, &keys)
	if err != nil {
		t.Errorf("JSON unmarshal error: %v", err)
	}
	// 验证获取的keys是否正确
	expectedKeys := handler.GetListKeys()
	if len(keys) != len(expectedKeys) {
		t.Errorf("unexpected number of keys: expected %d, got %d", len(expectedKeys), len(keys))
	}

	for i := 0; i < len(keys); i++ {
		strA := string(keys[i])
		strB := string(expectedKeys[i])
		if strA != strB {
			t.Errorf("Get error: expected ok, got %s", string(strA))
		} else {
			value, err2 := handler.Get(keys[i])
			if err2 != nil {
				return
			}
			t.Logf("the key of string:%s------the value of the key:%s", string(strA), string(value))
		}
	}

}
