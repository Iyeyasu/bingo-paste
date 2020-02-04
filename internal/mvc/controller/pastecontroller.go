package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"bingo/internal/config"
	"bingo/internal/http/httpext"
	"bingo/internal/mvc/model"
	"bingo/internal/mvc/model/store"
	"bingo/internal/mvc/view"
	"bingo/internal/session"
)

// PasteController handles creating and displaying pastes.
type PasteController struct {
	store *store.PasteStore
	view  *view.PasteView
}

// NewPasteController creates a new PasteController.
func NewPasteController(store *store.PasteStore) *PasteController {
	ctrl := new(PasteController)
	ctrl.store = store
	ctrl.view = view.NewPasteView()
	return ctrl
}

// ServeWritePage serves the page for creating pastes.
func (ctrl *PasteController) ServeWritePage(w http.ResponseWriter, r *http.Request) {
	// Because Viewer role can't write pastes, redirect to a more meaningful page
	user := session.User(r)
	if user != nil && user.Role == config.RoleViewer {
		http.Redirect(w, r, "/pastes", http.StatusFound)
		return
	}

	ctx := ctrl.view.NewWritePasteContext(r)
	ctrl.view.Write.Render(w, ctx)
}

// ServeViewPage serves the page for viewing individual pastes.
func (ctrl *PasteController) ServeViewPage(w http.ResponseWriter, r *http.Request) {
	id, err := httpext.ParseID(r)
	if err != nil {
		httpext.InternalError(w, err.Error())
		return
	}

	paste, err := ctrl.store.FindByID(id)
	if err != nil {
		httpext.InternalError(w, err.Error())
		return
	}

	ctx := ctrl.view.NewViewPasteContext(r, paste)
	ctrl.view.View.Render(w, ctx)
}

// ServeRawPaste serves the raw text content of individual pastes.
func (ctrl *PasteController) ServeRawPaste(w http.ResponseWriter, r *http.Request) {
	id, err := httpext.ParseID(r)
	if err != nil {
		httpext.InternalError(w, err.Error())
		return
	}

	paste, err := ctrl.store.FindByID(id)
	if err != nil {
		httpext.InternalError(w, err.Error())
		return
	}

	httpext.WriteText(w, []byte(paste.RawContent))
}

// ServeListPage serves the page for viewing a list of pastes.
func (ctrl *PasteController) ServeListPage(w http.ResponseWriter, r *http.Request) {
	limit, offset := httpext.ParseRange(r)
	filter := httpext.ParseFilter(r)

	var pastes []*model.Paste
	var err error
	if filter == "" {
		pastes, err = ctrl.store.FindRange(limit, offset)
	} else {
		pastes, err = ctrl.store.Search(filter, limit, offset)
	}

	if err != nil {
		httpext.InternalError(w, err.Error())
	}

	ctx := ctrl.view.NewListPastesContext(r, pastes)
	ctrl.view.List.Render(w, ctx)
}

// CreatePaste creates a new paste.
func (ctrl *PasteController) CreatePaste(w http.ResponseWriter, r *http.Request) {
	template, err := parseTemplate(r)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln(err))
		return
	}

	paste, err := ctrl.store.Insert(template)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln(err))
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/pastes/%d", paste.ID), 303)
}

func parseTemplate(r *http.Request) (*model.PasteTemplate, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse paste: %s", err)
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse paste: %s", err)
	}

	duration, err := strconv.ParseInt(values.Get("expiry"), 10, 64)
	if err != nil {
		duration = 0
	}

	visibility, err := strconv.Atoi(values.Get("visibility"))
	if err != nil {
		visibility = int(config.VisibilityUnlisted)
	}

	pasteTmpl := model.PasteTemplate{
		Title:      values.Get("title"),
		RawContent: values.Get("content"),
		Visibility: config.Visibility(visibility),
		Duration:   time.Duration(duration),
		Language:   values.Get("language"),
	}
	return &pasteTmpl, nil
}
