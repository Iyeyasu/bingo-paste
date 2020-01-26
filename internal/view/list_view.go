package view

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Iyeyasu/bingo-paste/internal/model"
	http_util "github.com/Iyeyasu/bingo-paste/internal/util/http"
	template_util "github.com/Iyeyasu/bingo-paste/internal/util/template"
	"github.com/julienschmidt/httprouter"
)

var (
	pastesPerPage int64 = 10
)

// ListView .
type ListView struct {
	store    *model.PasteStore
	template *template.Template
}

// NewListView .
func NewListView(store *model.PasteStore) *ListView {
	ctx := template_util.NewTemplateContext()
	view := new(ListView)
	view.store = store
	view.template = template_util.PrerenderTemplate("list", ctx)
	return view
}

// ServeList .
func (view *ListView) ServeList(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if pastes, err := view.getPastes(params); err == nil {
		http_util.WriteTemplate(w, view.template, pastes)
	} else {
		http_util.WriteError(w, fmt.Sprintf("Failed to retrieve paste: %s", err.Error()))
	}
}

func (view *ListView) getPastes(params httprouter.Params) ([]*model.Paste, error) {
	limitParam := params.ByName("limit")
	searchParam := params.ByName("search")
	if limitParam == "" {
		return view.store.SelectPublicSlice(0, pastesPerPage, searchParam)
	}

	limit, err := strconv.ParseInt(limitParam, 10, 64)
	if err != nil {
		return nil, err
	}

	return view.store.SelectPublicSlice(limit, pastesPerPage, searchParam)
}
