package view

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Iyeyasu/bingo-paste/internal/config"
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
	name     string
	store    *model.PasteStore
	template *template.Template
}

// ListTemplateContext .
type ListTemplateContext struct {
	template_util.TemplateContext
	Pastes []*model.Paste
}

// NewListView .
func NewListView(store *model.PasteStore) *ListView {
	view := new(ListView)
	view.name = "list"
	view.store = store
	view.template = template_util.GetTemplate(view.name)
	return view
}

// ServeList .
func (view *ListView) ServeList(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	query := req.URL.Query()
	if pastes, err := view.getPastes(query); err == nil {
		ctx := ListTemplateContext{
			Pastes: pastes,
			TemplateContext: template_util.TemplateContext{
				View:   view.name,
				Config: config.Get(),
				Filter: query.Get("search"),
			},
		}
		http_util.WriteTemplate(w, view.template, ctx)
	} else {
		http_util.WriteError(w, fmt.Sprintf("Failed to retrieve paste: %s", err.Error()))
	}
}

func (view *ListView) getPastes(query url.Values) ([]*model.Paste, error) {
	limitParam := query.Get("limit")
	searchParam := query.Get("search")
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
