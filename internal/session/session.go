package session

import (
	"errors"
	"net/http"
	"time"

	"bingo/internal/config"
	"bingo/internal/http/httpext"
	"bingo/internal/mvc/model"
	"bingo/internal/mvc/model/store"
	"bingo/internal/util/log"

	"github.com/alexedwards/scs/v2"
)

// Default is the default session.
var Default *Session

// Session handles user sessions.
type Session struct {
	Manager   *scs.SessionManager
	store     *RedisStore
	userStore *store.UserStore
}

// InitDefault initializes the default session.
func InitDefault(store *store.UserStore) {
	Default = New(store)
}

// New creates a new Session.
func New(store *store.UserStore) *Session {
	manager := scs.New()
	manager.Lifetime = 365 * 24 * time.Hour
	manager.Cookie.Name = config.Get().Authentication.Session.Name
	manager.Cookie.HttpOnly = true
	manager.Cookie.Path = "/"
	manager.Cookie.Persist = true
	manager.Cookie.SameSite = http.SameSiteLaxMode
	manager.Cookie.Secure = config.Get().Authentication.Session.SecureCookie
	manager.ErrorFunc = sessionError

	if config.Get().Authentication.Session.Store.Type == "redis" {
		manager.Store = NewRedisStore()
	}

	session := new(Session)
	session.Manager = manager
	session.userStore = store
	return session
}

// User returns the active user for the session.
func User(r *http.Request) *model.User {
	if Default == nil {
		return nil
	}

	cachedUser := GetRequestValue(r, "user")
	if cachedUser != nil {
		return cachedUser.(*model.User)
	}

	if !Default.Manager.Exists(r.Context(), "user_id") {
		return nil
	}

	id := Default.Manager.Get(r.Context(), "user_id").(int64)
	user, err := Default.userStore.FindByID(id)
	if err != nil {
		return nil
	}

	SetRequestValue(r, "user", user)
	return user
}

// Login sets the active user for the session.
func Login(r *http.Request, user *model.User) error {
	if Default == nil {
		return errors.New("failed to log in: no session configured")
	}

	err := renewToken(r)
	if err != nil {
		Default.Manager.Clear(r.Context())
		return err
	}

	Default.Manager.Put(r.Context(), "user_id", user.ID)
	return nil
}

// Logout clears the active user for the session.
func Logout(r *http.Request) error {
	if Default == nil {
		return errors.New("failed to log out: no session configured")
	}

	renewToken(r)
	Default.Manager.Clear(r.Context())
	return nil
}

// Make sure to renew token to prevent session fixation attack.
func renewToken(r *http.Request) error {
	err := Default.Manager.RenewToken(r.Context())
	if err != nil {
		log.Debugln("Failed to renew session token:", err)
	}
	return err
}

func sessionError(w http.ResponseWriter, r *http.Request, err error) {
	httpext.InternalError(w, err.Error())
}
