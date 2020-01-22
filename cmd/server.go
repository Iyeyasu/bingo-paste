package main

import (
	"log"
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/api"
	"github.com/Iyeyasu/bingo-paste/internal/model"
	"github.com/Iyeyasu/bingo-paste/internal/view"
	"github.com/julienschmidt/httprouter"
)

func main() {
	db, err := model.OpenDB()
	if err != nil {
		log.Fatalln(err)
	}

	router := httprouter.New()
	pasteStore := model.NewStore(db)

	pasteView := view.NewPasteView(router, pasteStore)
	pasteView.Handle("/")
	pasteView.Handle("/view/:id")

	pasteEndPoint := api.NewPasteEndPoint(router, pasteStore)
	pasteEndPoint.Handle("/api/v1/paste/")

	router.GET("/favicon.ico", faviconHandler)
	log.Fatal(http.ListenAndServe(":80", router))
}

func faviconHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	http.ServeFile(w, r, "web/img/favicon.ico")
}
