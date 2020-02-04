package session

import (
	"bingo/internal/config"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// RedisStore represents a persistant session store.
type RedisStore struct {
	client *redis.Client
	prefix string
}

// NewRedisStore returns a new RedisStore instance.
func NewRedisStore() *RedisStore {
	conf := config.Get().Authentication.Session.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.Database,
	})

	return &RedisStore{
		client: client,
		prefix: "scs:session:",
	}
}

// Find returns the data for a given session token from the RedisStore instance.
// If the session token is not found or is expired, the returned exists flag
// will be set to false.
func (store *RedisStore) Find(token string) (b []byte, exists bool, err error) {
	bytes, err := store.client.Get(store.prefix + token).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return bytes, true, nil
}

// Commit adds a session token and data to the RedisStore instance with the
// given expiry time. If the session token already exists then the data and
// expiry time are updated.
func (store *RedisStore) Commit(token string, b []byte, expiry time.Time) error {
	return store.client.Set(store.prefix+token, b, expiry.Sub(time.Now())).Err()
}

// Delete removes a session token and corresponding data from the RedisStore
// instance.
func (store *RedisStore) Delete(token string) error {
	return store.client.Del(store.prefix + token).Err()
}
