package http_handler

import (
	"encoding/json"
	"github.com/ByteStorage/FlyDB/engine"
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
