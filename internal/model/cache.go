package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var (
	connectionRetries = 45
	maxMemory         = "128mb"
	evictionPolicy    = "allkeys-lru"
	appendOnly        = "no"
)

// NewCache opens a new redis connection.
func NewCache() *redis.Client {

	if !config.Get().Cache.Enabled {
		return nil
	}

	addr := fmt.Sprintf("%s:%d", config.Get().Cache.Host, config.Get().Cache.Port)
	cache := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Get().Cache.Password,
		DB:       config.Get().Cache.Database,
	})

	log.Infof("Connecting to (addr = redis %s)", cache.Options().Addr)
	err := pollCache(cache)
	if err != nil {
		log.Fatalf("Failed to connect to cache: %s", err)
	}

	return cache
}

func pollCache(cache *redis.Client) error {
	log.Infof("Trying to connect to cache for %d seconds", connectionRetries)

	for i := 0; i <= connectionRetries; i++ {
		_, err := cache.Ping().Result()
		if err == nil {
			log.Info("Connected to redis")
			return nil
		}

		log.Info(err)
		time.Sleep(time.Second)
	}
	return errors.New("failed to connect to redis")
}

func configureCache(cache *redis.Client) {
	log.Debug("Configuring redis")

	log.Debugf("Setting maxmemory to %s", maxMemory)
	cache.ConfigSet("maxmemory", maxMemory)
	log.Debugf("Setting eviction policy to %d", evictionPolicy)
	cache.ConfigSet("maxmemory-policy", evictionPolicy)
	log.Debugf("Setting appendonly to %d", appendOnly)
	cache.ConfigSet("appendonly", appendOnly)
	log.Debugf("Setting save to ''")
	cache.ConfigSet("save", "")
}
