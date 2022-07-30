package pkg

import (
	"fmt"
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
	data, found := m.cache.Get(key)
	return data.(string), found
}

func (m *MemCache) Set(value string) string {

	randomID := m.generateRandomID()
	m.cache.Set(randomID, value, cache.DefaultExpiration)
	return randomID
}

func (m *MemCache) Delete(key string) {
	//TO-DO
}

func (m *MemCache) generateRandomID() string {
	return uuid.New().String()
}

func (m *MemCache) DumpAllItems() {
	items := m.cache.Items()

	for k, v := range items {
		fmt.Printf("Key : %s , Value : %v \n", k, v)
	}
}
