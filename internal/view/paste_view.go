package view

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/Iyeyasu/bingo-paste/internal/model"
	"github.com/julienschmidt/httprouter"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
)

// PasteView is the main view for editing and viewing pastes.
type PasteView struct {
	name           string
	router         *httprouter.Router
	store          *model.PasteStore
	editorTemplate []byte
	viewerTemplate *template.Template
}

// NewPasteView creates a new view for pastes.
func NewPasteView(router *httprouter.Router, store *model.PasteStore) *PasteView {
	log.Printf("Creating Paste view")

	view := new(PasteView)
	view.name = "Paste view"
	view.router = router
	view.store = store
	view.editorTemplate = []byte(view.renderTemplate(false))
	view.viewerTemplate = template.Must(template.New("").Parse(view.renderTemplate(true)))

	return view
}

// Handle sets up the HTTP request handling for the given URL.
func (view *PasteView) Handle(url string) {
	view.router.GET(url, view.render)
}

func (view *PasteView) render(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	if id != "" {
		view.renderViewer(w, id)
	} else {
		view.renderEditor(w)
	}
}

func (view *PasteView) renderViewer(w http.ResponseWriter, idStr string) {
	log.Printf("Rendering viewer '%s'.", view.name)

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paste, err := view.store.Select(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = view.viewerTemplate.Execute(w, paste)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (view *PasteView) renderEditor(w http.ResponseWriter) {
	log.Printf("Rendering view '%s'.", view.name)

	w.Write(view.editorTemplate)
}

func (view *PasteView) renderTemplate(readOnly bool) string {
	log.Printf("Rendering view template '%s'.", view.name)

	ctx := newRenderContext()
	ctx.ReadOnly = readOnly

	var buf bytes.Buffer
	tmpl := template.Must(template.ParseGlob("web/template/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("web/css/*.css"))
	tmpl.ExecuteTemplate(&buf, "index.html", ctx)
	tmplMini := view.minifyTemplate(buf.String())
	return tmplMini
}

func (view *PasteView) minifyTemplate(str string) string {
	log.Printf("Minifying view template '%s'.", view.name)

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

	res, err := m.String("text/html", str)

	if err != nil {
		log.Panicln("Failed to minify HTML template")
	}

	return res
}
