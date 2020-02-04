package session

import (
	"errors"
	"net/http"
	"time"

	"bingo/internal/config"
	"bingo/internal/mvc/model"
	"bingo/internal/mvc/model/store"
	"bingo/internal/util/log"

	"github.com/alexedwards/scs/v2"
)

var session *Session

// Session handles user sessions.
type Session struct {
	Manager   *scs.SessionManager
	userStore *store.UserStore
}

// Init initializes the default session.
func Init(store *store.UserStore) {
	session = New(store)
}

// Get returns the default session.
func Get() *scs.SessionManager {
	return session.Manager
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

	if config.Get().Authentication.Session.Store == "redis" {
		manager.Store = NewRedisStore()
	} else if config.Get().Authentication.Session.Store == "db" {
		manager.Store = NewPostgresStore(store.Database.DB)
	}

	session := new(Session)
	session.Manager = manager
	session.userStore = store
	return session
}

// User returns the active user for the session.
func User(r *http.Request) *model.User {
	cachedUser := GetRequestValue(r, "user")
	if cachedUser != nil {
		return cachedUser.(*model.User)
	}

	if !session.Manager.Exists(r.Context(), "user_id") {
		return nil
	}

	id := session.Manager.Get(r.Context(), "user_id").(int64)
	user, err := session.userStore.FindByID(id)
	if err != nil {
		return nil
	}

	SetRequestValue(r, "user", user)
	return user
}

// Login sets the active user for the session.
func Login(r *http.Request, user *model.User) error {
	if session.Manager == nil {
		return errors.New("no session configured")
	}

	err := renewToken(r)
	if err != nil {
		session.Manager.Clear(r.Context())
		return err
	}

	session.Manager.Put(r.Context(), "user_id", user.ID)
	return nil
}

// Logout clears the active user for the session.
func Logout(r *http.Request) error {
	if session == nil {
		return errors.New("no session configured")
	}

	renewToken(r)
	session.Manager.Clear(r.Context())
	return nil
}

// Make sure to renew token to prevent session fixation attack.
func renewToken(r *http.Request) error {
	err := session.Manager.RenewToken(r.Context())
	if err != nil {
		log.Debugln("Failed to renew session token:", err)
	}
	return err
}
