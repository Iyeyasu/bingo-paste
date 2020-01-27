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
	ctx := template_util.NewTemplateContext("list")
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
		if searchParam == "" {
			return view.store.SelectList(0, pastesPerPage)
		}
		return view.store.SearchList(searchParam, 0, pastesPerPage)
	}

	limit, err := strconv.ParseInt(limitParam, 10, 64)
	if err != nil {
		return nil, err
	}

	if searchParam == "" {
		return view.store.SelectList(limit, pastesPerPage)
	}
	return view.store.SearchList(searchParam, limit, pastesPerPage)
}
