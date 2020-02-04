package httpext

import (
	"bingo/internal/mvc/model"
	"bingo/internal/session"
	"encoding/json"
	"html/template"
	"net/http"
)

// WriteDefaultHeaders writes default headers to the HTTP response.
func WriteDefaultHeaders(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}

// WriteRaw writes the specified content type to the HTTP response.
func WriteRaw(w http.ResponseWriter, contentType string, output []byte) error {
	WriteDefaultHeaders(w, contentType)

	_, err := w.Write(output)
	return err
}

// WriteText writes raw plain-text to the HTTP response.
func WriteText(w http.ResponseWriter, output []byte) error {
	return WriteRaw(w, "text/plain", output)
}

// WriteHTML writes raw HTML to the HTTP response.
func WriteHTML(w http.ResponseWriter, output []byte) error {
	return WriteRaw(w, "text/html", output)
}

// WriteJSON writes raw JSON to the HTTP response.
func WriteJSON(w http.ResponseWriter, output interface{}) error {
	WriteDefaultHeaders(w, "application/json")

	encode := json.NewEncoder(w)
	return encode.Encode(output)
}

// WriteTemplate writes the given template to the HTTP response.
func WriteTemplate(w http.ResponseWriter, tmpl *template.Template, ctx interface{}) error {
	WriteDefaultHeaders(w, "text/html")

	return tmpl.Execute(w, ctx)
}

// WriteErrorNotification reloads the current page and inserts an error notification.
func WriteErrorNotification(w http.ResponseWriter, r *http.Request, title string, content string) {
	notification := model.NewErrorNotification(title, content)
	session.Get().Put(r.Context(), model.NotificationKey, notification)
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}

// WriteSuccessNotification reloads the current page and inserts a success notification.
func WriteSuccessNotification(w http.ResponseWriter, r *http.Request, title string, content string) {
	notification := model.NewSuccessNotification(title, content)
	session.Get().Put(r.Context(), model.NotificationKey, notification)
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}

// WriteNotification reloads the current page and inserts a notification.
func WriteNotification(w http.ResponseWriter, r *http.Request, notification *model.Notification) {
	session.Get().Put(r.Context(), model.NotificationKey, notification)
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}
