package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Iyeyasu/bingo-paste/internal/model"
	"github.com/julienschmidt/httprouter"
)

// PasteEndPoint represents a REST API endpoint for retrieving pastes.
type PasteEndPoint struct {
	router *httprouter.Router
	store  *model.PasteStore
}

// NewPasteEndPoint creates a new REST API endpoint for pastes.
func NewPasteEndPoint(router *httprouter.Router, store *model.PasteStore) *PasteEndPoint {
	endPoint := new(PasteEndPoint)
	endPoint.router = router
	endPoint.store = store
	return endPoint
}

// Handle sets the URI for the end point.
func (endPoint *PasteEndPoint) Handle(uri string) {
	endPoint.router.GET(uri, endPoint.get)
	endPoint.router.POST(uri, endPoint.post)
}

func (endPoint *PasteEndPoint) get(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("Retrieving paste...")

	id := int64(0)
	paste, err := endPoint.store.Select(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encode := json.NewEncoder(w)
	err = encode.Encode(paste)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Paste retrieved.")
}

func (endPoint *PasteEndPoint) post(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("Creating paste...")

	mime := req.Header.Get("Content-Type")
	log.Printf("Decoding request (Content-Type: %s)...", mime)

	var err error
	var paste *model.Paste
	switch mime {
	case "application/x-www-form-urlencoded":
		paste, err = decodeURL(req)
	default:
		err = fmt.Errorf("unrecognized Content-Type '%s'", mime)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := endPoint.store.Insert(paste)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("/view/%d", id)
	log.Printf("Paste created. Redirecting to %s", url)
	http.Redirect(w, req, url, 303)
}

func decodeJSON(req *http.Request) (*model.Paste, error) {
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	paste := new(model.Paste)
	err := decoder.Decode(paste)
	if err != nil {
		return nil, err
	}

	if decoder.More() {
		return nil, errors.New("extraneous data after JSON object")
	}

	return paste, nil
}

func decodeURL(req *http.Request) (*model.Paste, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	decoded, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	lifetime, err := strconv.ParseInt(decoded.Get("retention"), 10, 64)
	if err != nil {
		return nil, err
	}

	paste := model.Paste{
		Title:           decoded.Get("title"),
		Content:         decoded.Get("content"),
		IsPublic:        decoded.Get("visibility") == "public",
		LifetimeSeconds: int64(lifetime),
		Syntax:          decoded.Get("syntax"),
	}
	return &paste, nil
}

// func decodeID(req *http.Request) {
// 	re, _ := regexp.Compile("/id/(.*)")
// 	values := re.FindStringSubmatch(req.)
// 	if len(values) < 2 {

// 	}
// }
