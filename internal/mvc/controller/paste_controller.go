package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Iyeyasu/bingo-paste/internal/http/httpext"
	"github.com/Iyeyasu/bingo-paste/internal/mvc/model"
	"github.com/Iyeyasu/bingo-paste/internal/mvc/view"
)

// PasteController handles creating and displaying pastes.
type PasteController struct {
	store *model.PasteStore
	view  *view.PasteView
}

// NewPasteController creates a new PasteController.
func NewPasteController(store *model.PasteStore) *PasteController {
	ctrl := new(PasteController)
	ctrl.store = store
	ctrl.view = view.NewPasteView()
	return ctrl
}

// ServeEditPage serves the page for creating pastes.
func (ctrl *PasteController) ServeEditPage(w http.ResponseWriter, r *http.Request) {
	ctx := ctrl.view.NewPasteEditorContext()
	ctrl.view.Edit.Render(w, ctx)
}

// ServeViewPage serves the page for viewing individual pastes.
func (ctrl *PasteController) ServeViewPage(w http.ResponseWriter, r *http.Request) {
	id, err := httpext.ParseID(r)
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to read paste:", err.Error()))
		return
	}

	paste, err := ctrl.store.FindByID(id)
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to read paste:", err.Error()))
		return
	}

	ctx := ctrl.view.NewPasteViewerContext(paste)
	ctrl.view.View.Render(w, ctx)
}

// ServeRawPaste serves the raw text content of individual pastes.
func (ctrl *PasteController) ServeRawPaste(w http.ResponseWriter, r *http.Request) {
	id, err := httpext.ParseID(r)
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to read raw paste:", err.Error()))
		return
	}

	paste, err := ctrl.store.FindByID(id)
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to read raw paste:", err.Error()))
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
		httpext.WriteError(w, fmt.Sprintf("failed to read pastes: %s", err.Error()))
	}

	ctx := ctrl.view.NewPasteListContext(pastes, filter)
	ctrl.view.List.Render(w, ctx)
}

// CreatePaste creates a new paste.
func (ctrl *PasteController) CreatePaste(w http.ResponseWriter, r *http.Request) {
	template, err := parseTemplate(r)
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to create new paste:", err))
		return
	}

	paste, err := ctrl.store.Insert(template)
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to create new paste:", err))
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
		return nil, fmt.Errorf("failed to parse paste: %s", err)
	}

	template := model.PasteTemplate{
		Title:      values.Get("title"),
		RawContent: values.Get("content"),
		IsPublic:   values.Get("visibility") == "public",
		Duration:   time.Duration(duration),
		Language:   values.Get("language"),
	}
	return &template, nil
}
