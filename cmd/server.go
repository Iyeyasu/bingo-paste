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

	log.Fatal(http.ListenAndServe(":80", nil))
}
