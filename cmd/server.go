package main

import (
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/api"
	"github.com/Iyeyasu/bingo-paste/internal/model"
	"github.com/Iyeyasu/bingo-paste/internal/view"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func main() {
	// log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	database := model.NewDatabase()
	router := httprouter.New()
	pasteView := view.NewPasteView(database.Pastes)
	errorView := view.NewErrorView()

	router.GET("/favicon.ico", faviconHandler)
	router.GET("/", pasteView.ServeEditor)
	router.GET("/view/:id", pasteView.ServePaste)
	router.GET("/view/:id/raw", pasteView.ServeRawPaste)
	router.NotFound = errorView

	pasteEndPoint := api.NewPasteEndPoint(router, database.Pastes)
	pasteEndPoint.Handle("/api/v1/paste/")

	log.Fatal(http.ListenAndServe(":80", router))
}

func faviconHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	http.ServeFile(w, r, "web/img/favicon.ico")
}
