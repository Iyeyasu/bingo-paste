package session

import (
	"net/http"
	"time"

	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/http/httpext"
	"github.com/Iyeyasu/bingo-paste/internal/mvc/model"
	"github.com/Iyeyasu/bingo-paste/internal/mvc/model/store"
	"github.com/Iyeyasu/bingo-paste/internal/util/log"
	"github.com/alexedwards/scs/v2"
)

// Default is the default session.
var Default *Session

// Session handles user sessions.
type Session struct {
	Manager *scs.SessionManager
	store   *store.UserStore
}

// InitDefault initializes the default session.
func InitDefault(store *store.UserStore) {
	Default = New(store)
}

// New creates a new Session.
func New(store *store.UserStore) *Session {
	manager := scs.New()
	manager.Lifetime = 24 * time.Hour
	manager.Cookie.Name = config.Get().Authentication.SessionName
	manager.Cookie.HttpOnly = true
	manager.Cookie.Path = "/"
	manager.Cookie.Persist = true
	manager.Cookie.SameSite = http.SameSiteLaxMode
	manager.Cookie.Secure = config.Get().Authentication.SecureCookie
	manager.ErrorFunc = sessionError

	session := new(Session)
	session.Manager = manager
	session.store = store
	return session
}

// User returns the active user for the session.
func User(r *http.Request) *model.User {
	cachedUser := GetRequestValue(r, "user")
	if cachedUser != nil {
		return cachedUser.(*model.User)
	}

	if !Default.Manager.Exists(r.Context(), "user_id") {
		return nil
	}

	id := Default.Manager.Get(r.Context(), "user_id").(int64)
	user, err := Default.store.FindByID(id)
	if err != nil {
		return nil
	}

	SetRequestValue(r, "user", user)
	return user
}

// Login sets the active user for the session.
func Login(r *http.Request, user *model.User) error {
	// Make sure to renew token to prevent session fixation attack.
	err := Default.Manager.RenewToken(r.Context())
	if err != nil {
		log.Debugf("Failed to renew session token on login", err)
		Default.Manager.Clear(r.Context())
		return err
	}

	Default.Manager.Put(r.Context(), "user_id", user.ID)

	_, _, err = Default.Manager.Commit(r.Context())
	if err != nil {
		log.Debugf("Failed to commit login change", err)
		Default.Manager.Clear(r.Context())
		return err
	}

	return nil
}

// Logout clears the active user for the session.
func Logout(r *http.Request) {
	err := Default.Manager.RenewToken(r.Context())
	if err != nil {
		log.Debugf("Failed to renew session token on logout", err)
	}

	Default.Manager.Clear(r.Context())
}

func sessionError(w http.ResponseWriter, r *http.Request, err error) {
	httpext.InternalError(w, err.Error())
}
