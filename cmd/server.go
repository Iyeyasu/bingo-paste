package main

import (
	"log"
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/api"
	"github.com/Iyeyasu/bingo-paste/internal/view"
)

func main() {
	var pasteView view.PasteView
	pasteView.Handle("/")

	var copyView view.CopyView
	copyView.Handle("/view/")

	var endPoint api.PasteEndPoint
	endPoint.Handle("/api/v1/paste/")

	http.HandleFunc("/favicon.ico", faviconHandler)

	log.Fatal(http.ListenAndServe(":80", nil))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/img/favicon.ico")
}
