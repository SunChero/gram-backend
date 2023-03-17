package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
)

func (h *Handler) listStreams(w http.ResponseWriter, r *http.Request) {
	str := h.Svc.ListStreams()
	a, err := json.Marshal(str)
	if err != nil {
		return
	}
	_, err = w.Write(a)
	if err != nil {
		log.Println("cant process / request")
	}
}
func (h *Handler) updateStreams(w http.ResponseWriter, r *http.Request) {
	var s nats.StreamConfig
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		log.Print(`cant decode nats config to json`)
	}
	str, err := h.Svc.JS.UpdateStream(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	a, err := json.Marshal(str)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_, err = w.Write(a)
	if err != nil {
		log.Println("cant process / request")
	}

}
func (h *Handler) deleteStreams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str := vars["id"]
	_, err := h.Svc.DeleteStream(str)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
func (h *Handler) createStream(w http.ResponseWriter, r *http.Request) {
	var s nats.StreamConfig
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		log.Print(`cant decode nats config to json`)
	}
	str, err := h.Svc.CreateStream(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	a, err := json.Marshal(str)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_, err = w.Write(a)
	if err != nil {
		log.Println("cant process / request")
	}
}

func (h *Handler) joinStream(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	stream := v["stream"]
	consumer := r.Context().Value("name").(string)
	b, err := h.Svc.JoinStream(stream, consumer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	a, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Write([]byte(a))
}

func (h *Handler) leaveStream(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	stream := v["stream"]
	consumer := r.Context().Value("name").(string)
	err := h.Svc.LeaveStream(stream, consumer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
