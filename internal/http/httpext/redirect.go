package httpext

import (
	"bingo/internal/mvc/model"
	"bingo/internal/session"
	"net/http"
)

// RedirectWithNotify redirects to the given url and shows a notification on page load.
func RedirectWithNotify(w http.ResponseWriter, r *http.Request, url string, code int, note *model.Notification) {
	session.Get().Put(r.Context(), model.NotificationKey, note)
	http.Redirect(w, r, url, code)
}

// ReloadWithError reloads the current page and inserts an error notification.
func ReloadWithError(w http.ResponseWriter, r *http.Request, title string, content string) {
	notification := model.NewErrorNotification(title, content)
	session.Get().Put(r.Context(), model.NotificationKey, notification)
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}

// ReloadWithSuccess reloads the current page and inserts a success notification.
func ReloadWithSuccess(w http.ResponseWriter, r *http.Request, title string, content string) {
	notification := model.NewSuccessNotification(title, content)
	session.Get().Put(r.Context(), model.NotificationKey, notification)
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}

// ReloadWithNotification reloads the current page and inserts a notification.
func ReloadWithNotification(w http.ResponseWriter, r *http.Request, notification *model.Notification) {
	session.Get().Put(r.Context(), model.NotificationKey, notification)
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}
