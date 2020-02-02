package view

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/Iyeyasu/bingo-paste/internal/http/httpext"
	"github.com/Iyeyasu/bingo-paste/internal/util/fmtutil"
	"github.com/Iyeyasu/bingo-paste/internal/util/log"
)

// Page represents a single HTML template page.
type Page struct {
	Name     string
	Template *template.Template
}

// NewPage creates a new Page.
func NewPage(name string, paths []string) *Page {
	page := new(Page)
	page.Name = name
	page.Template = newTemplate(paths)
	return page
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
	log.Tracef("Template Function: unescaped HTML templates", str)
	return template.HTML(str)
}

func duration(duration time.Duration) int64 {
	result := int64(duration)
	log.Tracef("Template Function: formatted duration %s to %d", duration, result)
	return result
}

func formatExpiry(duration time.Duration, limit int) string {
	var result string
	if duration <= 0 {
		result = "Read Once"
	} else {
		result = fmtutil.FormatDuration(duration, limit)
	}

	log.Tracef("Template Function: formatted expiry duration %d -> %s", duration, result)
	return result
}

func formatPastDate(date time.Time) string {
	result := fmt.Sprintf("%s ago", fmtutil.FormatDuration(time.Now().Sub(date), 1))
	log.Tracef("Template Function: formatted past date %s", result)
	return result
}
