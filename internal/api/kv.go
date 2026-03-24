package api

import (
	"encoding/json"
	"net/http"

	"github.com/Saad7890-web/meridian/internal/storage"
)

type KVHandler struct {
	Store *storage.Store
}

func NewKVHandler(store *storage.Store) *KVHandler {
	return &KVHandler{Store: store}
}


func (h *KVHandler) Put(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	if key == "" {
		http.Error(w, "key required", http.StatusBadRequest)
		return
	}

	err := h.Store.Put(key, value)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte("OK"))
}


func (h *KVHandler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	entry, ok := h.Store.Get(key)
	if !ok {
		http.Error(w, "not found", 404)
		return
	}

	json.NewEncoder(w).Encode(entry)
}


func (h *KVHandler) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	err := h.Store.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte("deleted"))
}