package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/model"
)

// PasteEndPoint represents a REST API endpoint for retrieving pastes.
type PasteEndPoint struct {
}

// Handle sets the end point to handle REST requests to the passed URI.
func (endPoint *PasteEndPoint) Handle(uri string) {
	http.HandleFunc(uri, handle)
}

func handle(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		get(w, req)
	case "POST":
		post(w, req)
	case "DELETE":
		delete(w, req)
	default:
		notFound(w, req)
	}
}

func get(w http.ResponseWriter, req *http.Request) {
	var paste model.Paste
	encode := json.NewEncoder(w)
	encode.Encode(paste)
}

func post(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	var paste model.Paste
	err := decoder.Decode(&paste)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if decoder.More() {
		http.Error(w, "extraneous data after JSON object", http.StatusBadRequest)
		return
	}

	log.Println(paste)
}

func delete(w http.ResponseWriter, req *http.Request) {
}

func notFound(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not found"}`))
}
