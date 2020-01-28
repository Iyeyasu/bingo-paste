package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/model"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// PasteEndPoint represents a REST API endpoint for retrieving pastes.
type PasteEndPoint struct {
	store        *model.PasteStore
	uri          string
	defaultLimit int64
}

// NewPasteEndPoint creates a new REST API endpoint for pastes.
func NewPasteEndPoint(store *model.PasteStore) *PasteEndPoint {
	endPoint := new(PasteEndPoint)
	endPoint.uri = path.Join(config.Get().API.URI, endPoint.BaseURI())
	endPoint.store = store
	endPoint.defaultLimit = 10
	log.Debugf("Creating paste end point (%s)", endPoint.URI())
	return endPoint
}

// URI returns the URI of the end point.
func (endPoint *PasteEndPoint) URI() string {
	return endPoint.uri
}

// BaseURI returns the base URI of the end point.
func (endPoint *PasteEndPoint) BaseURI() string {
	return "/pastes"
}

// CreatePaste creates a new paste.
func (endPoint *PasteEndPoint) CreatePaste(r *http.Request) (*model.Paste, error) {
	log.Debugf("Creating a new paste using endpoint %s", endPoint.URI())

	paste, err := decodePaste(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create new paste: %s", err)
	}

	return endPoint.store.Insert(paste)
}

// GetPaste returns a single paste.
func (endPoint *PasteEndPoint) GetPaste(r *http.Request) (*model.Paste, error) {
	log.Debugf("Retrieving a paste from endpoint %s", endPoint.URI())

	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve paste: %s", err)
	}

	return endPoint.store.Select(id)
}

// GetPastes returns a range of pastes.
func (endPoint *PasteEndPoint) GetPastes(r *http.Request) ([]*model.Paste, error) {
	log.Debugf("Retrieving pastes using endpoint %s", endPoint.URI())

	query := r.URL.Query()
	limitParam := query.Get("limit")
	offsetParam := query.Get("offset")
	searchParam := query.Get("search")

	limit, err := strconv.ParseInt(limitParam, 10, 64)
	if err != nil {
		limit = endPoint.defaultLimit
	}

	offset, err := strconv.ParseInt(offsetParam, 10, 64)
	if err != nil {
		offset = 0
	}

	if searchParam == "" {
		return endPoint.store.SelectList(limit, offset)
	}
	return endPoint.store.SearchList(searchParam, limit, offset)
}

func decodePaste(r *http.Request) (*model.Paste, error) {
	mime := r.Header.Get("Content-Type")
	log.Debugf("Decoding paste using Content-Type: %s", mime)

	switch mime {
	case "application/x-www-form-urlencoded":
		return decodeURL(r)
	case "application/json":
		return decodeJSON(r)
	default:
		return nil, fmt.Errorf("failed to decode paste: unrecognized Content-Type '%s'", mime)
	}
}

func decodeJSON(r *http.Request) (*model.Paste, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	paste := new(model.Paste)
	err := decoder.Decode(paste)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %s", err)
	}

	if decoder.More() {
		return nil, errors.New("failed to decode JSON: extraneous data after JSON object")
	}

	return paste, nil
}

func decodeURL(r *http.Request) (*model.Paste, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode URL: %s", err)
	}

	decoded, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to decode URL: %s", err)
	}

	duration, err := strconv.ParseInt(decoded.Get("expiry"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode URL: %s", err)
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
