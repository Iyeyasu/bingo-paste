package main

import (
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/api"
	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/model"
	"github.com/Iyeyasu/bingo-paste/internal/view"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(config.Get().LogLevel)
	if config.Get().LogLevel != log.InfoLevel {
		log.SetReportCaller(true)
	}

	db := model.NewDatabase()
	pasteStore := model.NewPasteStore(db)

	router := httprouter.New()
	editorView := view.NewEditorView(pasteStore)
	viewerView := view.NewViewerView(pasteStore)
	listView := view.NewListView(pasteStore)
	errorView := view.NewErrorView()

	router.GET("/favicon.ico", faviconHandler)
	router.GET("/", editorView.ServeEditor)
	router.GET("/view/:id", viewerView.ServePaste)
	router.GET("/view/:id/raw", viewerView.ServeRawPaste)
	router.GET("/list", listView.ServeList)
	router.NotFound = errorView

	pasteEndPoint := api.NewPasteEndPoint(router, pasteStore)
	pasteEndPoint.Handle("/api/v1/paste/")

	log.Fatal(http.ListenAndServe(":80", router))
}

func faviconHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	http.ServeFile(w, r, "web/img/favicon.ico")
}
