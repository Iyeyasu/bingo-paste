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

// RenderTemplate renders the template with the given name, converting it to a string.
func RenderTemplate(name string, ctx interface{}) string {
	log.Debugf("Rendering template '%s'", name)

	tmpl := parseTemplate(name, newTemplate())
	builder := new(strings.Builder)
	err := tmpl.Execute(builder, ctx)
	if err != nil {
		log.Fatalf("Failed to render template '%s': %s", name, err.Error())
	}

	return html_util.Minify(builder.String())
}

// PrerenderTemplate renders the template with the given name, converting the
// resulting string into an new template. Useful when you have some parts of
// the page that are static and don't want to render them every time.
func PrerenderTemplate(name string, ctx interface{}) *template.Template {
	log.Debugf("Prerendering template '%s'", name)

	render := RenderTemplate(name, ctx)
	return template.Must(newTemplate().Parse(render))
}

func newTemplate() *template.Template {
	return template.New("index").Funcs(newFuncMap())
}

func parseTemplate(name string, tmpl *template.Template) *template.Template {
	tmpl = template.Must(tmpl.ParseGlob("web/css/*.css"))
	tmpl = template.Must(tmpl.ParseGlob("web/template/*.go.html"))
	tmpl = template.Must(tmpl.ParseGlob(fmt.Sprintf("web/template/%s/*.go.html", name)))
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
