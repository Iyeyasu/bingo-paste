package httpext

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/util/log"
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
		InternalError(w, fmt.Sprintf("Failed to write raw bytes to HTTP response: %s", err))
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

// InternalError writes an internal server error to the HTTP response.
func InternalError(w http.ResponseWriter, msg string) {
	log.Info(msg)
	http.Error(w, msg, http.StatusInternalServerError)
}

// UnauthorizedError writes an unauthorized error to the HTTP response.
func UnauthorizedError(w http.ResponseWriter) {
	msg := "You do not have the permission to view this page"
	log.Info(msg)
	http.Error(w, msg, http.StatusUnauthorized)
}

// WriteJSON writes raw JSON to the HTTP response.
func WriteJSON(w http.ResponseWriter, output interface{}) {
	WriteDefaultHeaders(w, "application/json")

	encode := json.NewEncoder(w)
	err := encode.Encode(output)
	if err != nil {
		InternalError(w, fmt.Sprintf("Failed to write JSON to HTTP response: %s", err))
	}
}

// WriteTemplate writes the given template to the HTTP response.
func WriteTemplate(w http.ResponseWriter, tmpl *template.Template, ctx interface{}) {
	WriteDefaultHeaders(w, "text/html")

	err := tmpl.Execute(w, ctx)
	if err != nil {
		InternalError(w, fmt.Sprintf("Failed to write template %s to HTTP response: %s", tmpl.Name(), err))
	}
}
