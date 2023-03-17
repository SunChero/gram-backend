package handler

import (
	"encoding/json"
	"net/http"
)

type Req struct {
	Stream   string `json:"stream"`
	Consumer string `json:"consumer"`
}

func (h *Handler) createConsumer(w http.ResponseWriter, r *http.Request) {
	var c Req
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	b, err := h.Svc.CreateConsumer(c.Stream, c.Consumer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	a, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Write([]byte(a))

}
func (h *Handler) deleteConsumer(w http.ResponseWriter, r *http.Request) {
	var c Req
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = h.Svc.DeleteConsumer(c.Stream, c.Consumer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

}
