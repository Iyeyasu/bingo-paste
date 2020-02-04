package view

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"bingo/internal/config"
	"bingo/internal/http/httpext"
	"bingo/internal/mvc/model"
	"bingo/internal/session"
	"bingo/internal/util/fmtutil"
)

// Page represents a single HTML template page.
type Page struct {
	Name     string
	Template *template.Template
}

// PageContext represents a rendering context for a page template.
type PageContext struct {
	Page          *Page
	Config        *config.Config
	CurrentUser   *model.User
	SearchFilter  string
	ShowSearchbar bool
}

// NewPage creates a new Page.
func NewPage(name string, paths []string) *Page {
	page := new(Page)
	page.Name = name
	page.Template = newTemplate(paths)
	return page
}

// NewPageContext creates a new PageContext.
func NewPageContext(r *http.Request, page *Page) PageContext {
	return PageContext{
		Page:          page,
		Config:        config.Get(),
		CurrentUser:   session.User(r),
		SearchFilter:  r.URL.Query().Get("search"),
		ShowSearchbar: true,
	}
}

// Render renders the page as a HTTP response using the given rendering context.
func (page *Page) Render(w http.ResponseWriter, ctx interface{}) {
	httpext.WriteTemplate(w, page.Template, ctx)
}

func newTemplate(paths []string) *template.Template {
	tmpl := template.New("index").Funcs(newFuncMap())
	for _, path := range paths {
		tmpl = template.Must(tmpl.ParseGlob(path))
	}
	return tmpl
}

func newFuncMap() template.FuncMap {
	return template.FuncMap{
		"duration":       duration,
		"formatExpiry":   formatExpiry,
		"formatPastDate": formatPastDate,
		"unescape":       unescape,
	}
}

func unescape(str string) template.HTML {
	return template.HTML(str)
}

func duration(duration time.Duration) int64 {
	result := int64(duration)
	return result
}

func formatExpiry(duration time.Duration, limit int) string {
	var result string
	if duration <= 0 {
		result = "Keep Forever"
	} else {
		result = fmtutil.FormatDuration(duration, limit)
	}

	return result
}

func formatPastDate(date time.Time) string {
	result := fmt.Sprintf("%s ago", fmtutil.FormatDuration(time.Now().Sub(date), 1))
	return result
}
