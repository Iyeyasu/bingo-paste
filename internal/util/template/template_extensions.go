package view

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	fmt_util "github.com/Iyeyasu/bingo-paste/internal/util/fmt"
	html_util "github.com/Iyeyasu/bingo-paste/internal/util/html"
	log "github.com/sirupsen/logrus"
)

// RenderTemplate renders the given template to a string.
func RenderTemplate(tmpl *template.Template, ctx interface{}) string {
	log.Debugf("Rendering template '%s'", tmpl.Name)

	builder := new(strings.Builder)
	err := tmpl.Execute(builder, ctx)
	if err != nil {
		log.Fatalf("Failed to render template '%s': %s", tmpl.Name, err.Error())
	}

	return html_util.Minify(builder.String())
}

// GetTemplate returns the template with the given name and globs.
func GetTemplate(name string, paths ...string) *template.Template {
	tmpl := newTemplate()
	for _, path := range paths {
		tmpl = template.Must(tmpl.ParseGlob(path))
	}
	return tmpl
}

func newTemplate() *template.Template {
	return template.New("index").Funcs(newFuncMap())
}

func newFuncMap() template.FuncMap {
	return template.FuncMap{
		"duration":       duration,
		"formatExpiry":   formatExpiry,
		"formatPastDate": formatPastDate,
		"unescape":       unescape,
	}
}

// Unescapes HTML so that we can inject syntax highlighted code to the viewer.
func unescape(str string) template.HTML {
	log.Tracef("Template Function: unescaped HTML templates", str)
	return template.HTML(str)
}

// Prevents duration from being rendered as a formatted string.
func duration(duration time.Duration) int64 {
	result := int64(duration)
	log.Tracef("Template Function: formatted duration %s to %d", duration, result)
	return result
}

// Formats expiry durations into a human readable form.
func formatExpiry(duration time.Duration, limit int) string {
	var result string
	if duration <= 0 {
		result = "Read Once"
	} else {
		result = fmt_util.FormatDuration(duration, limit)
	}

	log.Tracef("Template Function: formatted expiry duration %d -> %s", duration, result)
	return result
}

// Formats paste date into a human readable form of form "<time> ago".
func formatPastDate(date time.Time) string {
	result := fmt.Sprintf("%s ago", fmt_util.FormatDuration(time.Now().Sub(date), 1))
	log.Tracef("Template Function: formatted past date %s", result)
	return result
}
