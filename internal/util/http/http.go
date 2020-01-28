package util

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// WriteDefaultHeaders writes default headers to the HTTP response.
func WriteDefaultHeaders(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}

// WriteRaw writes the specified content type to the HTTP response.
func WriteRaw(w http.ResponseWriter, contentType string, output []byte) {
	WriteDefaultHeaders(w, contentType)
	_, err := w.Write(output)
	if err != nil {
		WriteError(w, fmt.Sprintf("Failed to write raw bytes to HTTP response: %s", err))
	}
}

// WriteJSON writes raw JSON to the HTTP response.
func WriteJSON(w http.ResponseWriter, output interface{}) {
	WriteDefaultHeaders(w, "application/json")
	encode := json.NewEncoder(w)
	err := encode.Encode(output)
	if err != nil {
		WriteError(w, fmt.Sprintf("Failed to write JSON to HTTP response: %s", err))
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
	log.Error(msg)
	http.Error(w, msg, http.StatusInternalServerError)
}

// WriteTemplate writes the given template to the HTTP response.
func WriteTemplate(w http.ResponseWriter, tmpl *template.Template, ctx interface{}) {
	WriteDefaultHeaders(w, "text/html")
	err := tmpl.Execute(w, ctx)
	if err != nil {
		WriteError(w, fmt.Sprintf("Failed to write template %s to HTTP response: %s", tmpl.Name(), err))
	}
}
