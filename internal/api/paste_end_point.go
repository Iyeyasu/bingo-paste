package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Iyeyasu/bingo-paste/internal/model"
	"github.com/Iyeyasu/bingo-paste/internal/util/http"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
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
	endPoint.router.GET(uri+":id", endPoint.get)
	endPoint.router.POST(uri, endPoint.post)
}

func (endPoint *PasteEndPoint) get(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	log.Debugf("GET paste request to %s", req.URL)

	id, err := util.ParseID(params)
	if err != nil {
		util.WriteError(w, fmt.Sprintf("Failed to GET paste from %s: %s", req.URL, err))
		return
	}

	paste, err := endPoint.store.Select(id)
	if err != nil {
		util.WriteError(w, fmt.Sprintf("Failed to GET paste %d from %s: %s", id, req.URL, err))
		return
	}

	log.Infof("GET paste %d from %s", req.URL)
	util.WriteJSON(w, paste)
}

func (endPoint *PasteEndPoint) post(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Debugf("POST paste request to %s", req.URL)

	paste, err := parsePaste(req)
	if err != nil {
		util.WriteError(w, fmt.Sprintf("Failed to POST paste to %s: %s", req.URL, err))
		return
	}

	id, err := endPoint.store.Insert(paste)
	if err != nil {
		util.WriteError(w, fmt.Sprintf("Failed to POST paste to %s: %s", req.URL, err))
		return
	}

	log.Infof("POST paste %d to %s", id, req.URL)
	util.Redirect(w, req, fmt.Sprintf("/view/%d", id), 303)
}

func parsePaste(req *http.Request) (*model.Paste, error) {
	log.Debug("Parsing POST paste request")

	mime := req.Header.Get("Content-Type")
	if mime == "" {
		return nil, errors.New("no Content-Type detected")
	}

	log.Debugf("Content-Type detected: %s", mime)
	switch mime {
	case "application/x-www-form-urlencoded":
		return decodeURL(req)
	case "application/json":
		return decodeJSON(req)
	default:
		return nil, fmt.Errorf("unrecognized Content-Type '%s'", mime)
	}
}

func decodeJSON(req *http.Request) (*model.Paste, error) {
	log.Debug("Parsing paste JSON data")

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
	log.Debug("Parsing paste URL encoded data")

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	decoded, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	duration, err := strconv.ParseInt(decoded.Get("expiry"), 10, 64)
	if err != nil {
		return nil, err
	}

	paste := model.Paste{
		Title:      decoded.Get("title"),
		RawContent: decoded.Get("content"),
		IsPublic:   decoded.Get("visibility") == "public",
		Duration:   time.Duration(duration),
		Language:   decoded.Get("language"),
	}
	return &paste, nil
}
