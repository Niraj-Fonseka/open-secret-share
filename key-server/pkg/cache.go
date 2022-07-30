package pkg

import (
	"time"

	"github.com/google/uuid"

	"github.com/patrickmn/go-cache"
)

type MemCache struct {
	cache *cache.Cache
}

func NewMemCache() *MemCache {

	//default expiration at 5 mins
	//eviction time 10 mins
	c := cache.New(1*time.Minute, 1*time.Minute)

	return &MemCache{
		cache: c,
	}
}

func (m *MemCache) Get(key string) (string, bool) {
	return m.cache.Get(key)
}

func (m *MemCache) Set(key, value string) string {

	randomID := m.generateRandomID()
	m.cache.Set(key, value, m.cache.DefaultExpiration)
	return randomID
}

func (m *MemCache) Delete(key string) {
	//TO-DO
}

func (m *MemCache) generateRandomID() string {
	return uuid.New().String()
}
