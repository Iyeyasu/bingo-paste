package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/data"
)

// PasteEndPoint represents a REST API endpoint for retrieving pastes.
type PasteEndPoint struct {
}

func (end PasteEndPoint) uri(prefix string) string {
	return prefix + "/paste"
}

func (end PasteEndPoint) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		get(rw, req)
	case "POST":
		post(rw, req)
	case "DELETE":
		delete(rw, req)
	default:
		notFound(rw, req)
	}
}

func get(rw http.ResponseWriter, req *http.Request) {
	var paste data.Paste
	encode := json.NewEncoder(rw)
	encode.Encode(paste)
}

func post(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	var paste data.Paste
	err := decoder.Decode(&paste)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if decoder.More() {
		http.Error(rw, "extraneous data after JSON object", http.StatusBadRequest)
		return
	}

	log.Println(paste)
}

func delete(rw http.ResponseWriter, req *http.Request) {
}

func notFound(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte(`{"message": "not found"}`))
}
