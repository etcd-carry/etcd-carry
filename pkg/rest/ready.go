package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Ready struct {
	sync.Mutex
	ready bool
}

func NewReady() *Ready {
	return &Ready{ready: false}
}

func (r *Ready) SetReady(isReady bool) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.ready = isReady
}

func (r *Ready) Handler(w http.ResponseWriter, _ *http.Request) {
	content, _ := json.Marshal(map[string]bool{"ready": r.ready})

	w.Header().Set("Content-Type", "application-json")
	if _, err := w.Write(content); err != nil {
		fmt.Println("response to ready request error:", err)
	}
}
