package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/model"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// UserEndPoint represents a REST API endpoint for retrieving users.
type UserEndPoint struct {
	store        *model.UserStore
	uri          string
	defaultLimit int64
}

// NewUserEndPoint creates a new REST API endpoint for users.
func NewUserEndPoint(store *model.UserStore) *UserEndPoint {
	endPoint := new(UserEndPoint)
	endPoint.uri = path.Join(config.Get().API.URI, endPoint.BaseURI())
	endPoint.store = store
	endPoint.defaultLimit = 10
	log.Debugf("Creating user end point (%s)", endPoint.URI())
	return endPoint
}

// URI returns the URI of the end point.
func (endPoint *UserEndPoint) URI() string {
	return endPoint.uri
}

// BaseURI returns the base URI of the end point.
func (endPoint *UserEndPoint) BaseURI() string {
	return "/users"
}

// CreateUser creates a new user.
func (endPoint *UserEndPoint) CreateUser(r *http.Request) (*model.User, error) {
	log.Debugf("Creating a new user using endpoint %s", endPoint.URI())

	user, err := endPoint.decodeUser(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create new user: %s", err)
	}

	return endPoint.store.Insert(user)
}

// ReadUser returns a single user.
func (endPoint *UserEndPoint) ReadUser(r *http.Request) (*model.User, error) {
	log.Debugf("Retrieving a user from endpoint %s", endPoint.URI())

	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to read user: %s", err)
	}

	return endPoint.store.Select(id)
}

// ReadUsers returns a range of users.
func (endPoint *UserEndPoint) ReadUsers(r *http.Request) ([]*model.User, error) {
	log.Debugf("Retrieving users using endpoint %s", endPoint.URI())

	query := r.URL.Query()
	limitParam := query.Get("limit")
	offsetParam := query.Get("offset")

	limit, err := strconv.ParseInt(limitParam, 10, 64)
	if err != nil {
		limit = endPoint.defaultLimit
	}

	offset, err := strconv.ParseInt(offsetParam, 10, 64)
	if err != nil {
		offset = 0
	}

	return endPoint.store.SelectMultiple(limit, offset)
}

// UpdateUser updates an existing user.
func (endPoint *UserEndPoint) UpdateUser(r *http.Request) (*model.User, error) {
	log.Debugf("Updating an existing user using endpoint %s", endPoint.URI())

	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %s", err)
	}

	user, err := endPoint.decodeUser(r)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %s", err)
	}
	user.ID = sql.NullInt64{Int64: id, Valid: true}

	return endPoint.store.Update(user)
}

// DeleteUser deletes an existing user.
func (endPoint *UserEndPoint) DeleteUser(r *http.Request) error {
	log.Debugf("Deleting an existing user using endpoint %s", endPoint.URI())

	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to delete user: %s", err)
	}

	return endPoint.store.Delete(id)
}

func (endPoint *UserEndPoint) decodeUser(r *http.Request) (*model.UserModel, error) {
	mime := r.Header.Get("Content-Type")
	log.Debugf("Decoding user using Content-Type: %s", mime)

	switch mime {
	case "application/x-www-form-urlencoded":
		return endPoint.decodeURL(r)
	case "application/json":
		return endPoint.decodeJSON(r)
	default:
		return nil, fmt.Errorf("failed to decode user: unrecognized Content-Type '%s'", mime)
	}
}

func (endPoint *UserEndPoint) decodeJSON(r *http.Request) (*model.UserModel, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	user := new(model.UserModel)
	err := decoder.Decode(user)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %s", err)
	}

	if decoder.More() {
		return nil, errors.New("failed to decode JSON: extraneous data after JSON object")
	}

	return user, nil
}

func (endPoint *UserEndPoint) decodeURL(r *http.Request) (*model.UserModel, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode URL: %s", err)
	}

	decoded, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to decode URL: %s", err)
	}

	user := model.UserModel{
		Password:       endPoint.decodeString(decoded.Get("password")),
		Name:           endPoint.decodeString(decoded.Get("name")),
		Mail:           endPoint.decodeString(decoded.Get("mail")),
		AuthExternalID: endPoint.decodeString(decoded.Get("auth_external_id")),
		AuthType:       endPoint.decodeInt32(decoded.Get("auth_type")),
		Role:           endPoint.decodeInt32(decoded.Get("role")),
		Theme:          endPoint.decodeInt32(decoded.Get("theme")),
	}

	return &user, nil
}

func (endPoint *UserEndPoint) decodeInt32(str string) sql.NullInt32 {
	val, err := strconv.Atoi(str)
	return sql.NullInt32{
		Int32: int32(val),
		Valid: err == nil,
	}
}

func (endPoint *UserEndPoint) decodeString(str string) sql.NullString {
	return sql.NullString{
		String: str,
		Valid:  str != "",
	}
}
