package http_handler

import (
	"encoding/json"
	"github.com/ByteStorage/FlyDB/engine"
	"io"
	"net/http"
)

type HttpHandler struct {
	*engine.DB
}

func NewHttpHandler(DB *engine.DB) *HttpHandler {
	return &HttpHandler{DB: DB}
}

// PutHandler 支持http进行Put
func (hs *HttpHandler) PutHandler(w http.ResponseWriter, r *http.Request) {
	type PutRequest struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	var putReq PutRequest
	err := json.NewDecoder(r.Body).Decode(&putReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if putReq.Key == "" {
		http.Error(w, "key is empty", http.StatusBadRequest)
		return
	}
	if putReq.Value == "" {
		http.Error(w, "value is empty", http.StatusBadRequest)
		return
	}
	err = hs.Put([]byte(putReq.Key), []byte(putReq.Value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte("ok"))
	if err != nil {
		return
	}
}

// DelHandler 支持http进行Delete
func (hs *HttpHandler) DelHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	if key == "" {
		http.Error(w, "key is empty", http.StatusBadRequest)
		return
	}
	err := hs.Delete([]byte(key))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte("ok"))
	if err != nil {
		return
	}
}

// GetHandler 支持http进行Get
func (hs *HttpHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	if key == "" {
		http.Error(w, "key is empty", http.StatusBadRequest)
		return
	}

	val, err := hs.Get([]byte(key))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(val)
	if err != nil {
		return
	}
}

// PostHandler 支持http进行Post
func (hs *HttpHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	type PostRequest struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	var postReq PostRequest
	err := json.NewDecoder(r.Body).Decode(&postReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)

	if postReq.Key == "" {
		http.Error(w, "key is empty", http.StatusBadRequest)
		return
	}
	if postReq.Value == "" {
		http.Error(w, "value is empty", http.StatusBadRequest)
		return
	}

	err = hs.Put([]byte(postReq.Key), []byte(postReq.Value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte("ok"))
	if err != nil {
		// 处理写入响应失败的错误
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetListKeysHandler 支持http获取数据库中所有键
func (hs *HttpHandler) GetListKeysHandler(w http.ResponseWriter, r *http.Request) {
	keys := hs.GetListKeys()
	if keys == nil {
		http.Error(w, "key is empty", http.StatusBadRequest)
		return
	}
	jsonKeys, err := json.Marshal(keys)
	if err != nil {
		// Handle the error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	_, err = w.Write(jsonKeys)
	if err != nil {
		// 处理写入响应失败的错误
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
