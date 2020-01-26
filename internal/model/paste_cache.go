package model

import (
	"fmt"
	"time"

	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var (
	cacheKey      = "paste:%d"
	cacheDuration = time.Duration(0) // Disable expiration
)

// PasteCache caches pastes for faster access.
type PasteCache struct {
	redis *redis.Client
}

// NewPasteCache creates a new PasteCache.
func NewPasteCache(redis *redis.Client) *PasteCache {
	cache := new(PasteCache)
	cache.redis = redis
	return cache
}

func (cache *PasteCache) set(paste *Paste) {
	if !config.Get().Cache.Enabled {
		return
	}

	log.Debugf("Storing paste %d in cache", paste.ID)

	key := fmt.Sprintf(cacheKey, paste.ID)
	result, err := cache.redis.SetNX(key, paste, cacheDuration).Result()
	if err != nil {
		log.Errorf("Failed to store paste %d in cache: %s", paste.ID, err)
	} else if !result {
		log.Debugf("Paste %d already cached", paste.ID)
	}
}

func (cache *PasteCache) get(id int64) *Paste {
	if !config.Get().Cache.Enabled {
		return nil
	}

	log.Debugf("Retrieving paste %d from cache", id)

	key := fmt.Sprintf(cacheKey, id)
	paste := new(Paste)
	err := cache.redis.Get(key).Scan(paste)
	if err == redis.Nil {
		log.Debugf("Failed to retrieve paste %d from cache: %s", id, err)
		return nil
	} else if err != nil {
		log.Errorf("Failed to retrieve paste %d from cache: %s", id, err)
		return nil
	}

	log.Debugf("Retrieved paste %d from cache", id)
	return paste
}

func (cache *PasteCache) invalidate(id int64) {
	if !config.Get().Cache.Enabled {
		return
	}

	log.Debugf("Invalidating paste %d in cache", id)

	key := fmt.Sprintf(cacheKey, id)
	err := cache.redis.Del(key).Err()
	if err != nil {
		log.Errorf("Failed to invalidate paste %d in cache: %s", id, err)
	}
}

func (cache *PasteCache) refresh(id int64) {
	if !config.Get().Cache.Enabled || cacheDuration == 0 {
		return
	}

	log.Debugf("Refreshing paste %d in cache", id)

	key := fmt.Sprintf(cacheKey, id)
	err := cache.redis.Expire(key, cacheDuration).Err()
	if err != nil {
		log.Errorf("Failed to refresh paste %d in cache: %s", id, err)
	}
}
