package util

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// Redirect redirects to the given URL.
func Redirect(w http.ResponseWriter, req *http.Request, url string, code int) {
	log.Debugf("Redirecting from %s to %s", req.URL, url)
	http.Redirect(w, req, url, code)
}

// ParseID parses ID from the given router params.
func ParseID(params httprouter.Params) (int64, error) {
	log.Debug("Parsing ID")
	return strconv.ParseInt(params.ByName("id"), 10, 64)
}

// WriteDefaultHeaders writes default headers to the HTTP response.
func WriteDefaultHeaders(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}

// WriteRaw writes the specified content type to the HTTP response.
func WriteRaw(w http.ResponseWriter, contentType string, output []byte) {
	WriteDefaultHeaders(w, contentType)
	_, err := w.Write(output)
	if err != nil {
		WriteError(w, "Failed to write raw bytes to HTTP response")
	}
}

// WriteJSON writes raw JSON to the HTTP response.
func WriteJSON(w http.ResponseWriter, output interface{}) {
	WriteDefaultHeaders(w, "application/json")
	encode := json.NewEncoder(w)
	err := encode.Encode(output)
	if err != nil {
		WriteError(w, "Failed to write JSON to HTTP response")
	}
}

// WriteText writes raw plain-text to the HTTP response.
func WriteText(w http.ResponseWriter, output []byte) {
	WriteRaw(w, "text/plain", output)
}

// WriteHTML writes raw HTML to the HTTP response.
func WriteHTML(w http.ResponseWriter, output []byte) {
	WriteRaw(w, "text/html", output)
}

// WriteError writes an error to the HTTP response.
func WriteError(w http.ResponseWriter, msg string) {
	log.Debug(msg)
	http.Error(w, msg, http.StatusNotFound)
}

// WriteTemplate writes the given template to the HTTP response.
func WriteTemplate(w http.ResponseWriter, tmpl *template.Template, ctx interface{}) {
	WriteDefaultHeaders(w, "text/html")
	err := tmpl.Execute(w, ctx)
	if err != nil {
		WriteError(w, fmt.Sprintf("Failed to write template %s to HTTP response", tmpl.Name()))
	}
}
