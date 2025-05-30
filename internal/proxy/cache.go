package proxy

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type Cache struct {
	Body      []byte
	Headers   http.Header
	Status    int
	CacheTime int64 // Unix timestamp for when the cache was created
}

var cache = make(map[string]*Cache)
var mu sync.RWMutex

// Get retrieves the cached response for a given URL.
func Get(url string) (*Cache, bool) {
	mu.RLock()
	defer mu.RUnlock()
	cachedResponse, found := cache[url]
	return cachedResponse, found
}

// Set caches the response for a given URL.
func Set(url string, response *Cache) {
	mu.Lock()
	defer mu.Unlock()
	cache[url] = response
}

// Clear removes the cached response for a given URL.
func Clear(url string) {
	mu.Lock()
	defer mu.Unlock()
	delete(cache, url)
}

// ClearAll removes all cached responses.
func ClearAll() {
	mu.Lock()
	defer mu.Unlock()
	log.Println("Clearing all cached responses", len(cache))
	cache = make(map[string]*Cache)
	log.Println("All cached responses cleared", len(cache))
}

// ClearOlderThan clears cached responses older than the 1-hour timestamp.
// This function is intended to be called periodically to maintain cache freshness.
// It removes entries that are older than 1 hour from the current time.
func ClearOlderThan() {
	log.Println("Starting periodic cache cleanup for entries older than 1 hour")
	go func() {
		for {
			mu.Lock()
			now := time.Now().Unix()
			hourAgo := now - 60*60 // 1 hour ago

			for key, cachedResponse := range cache {
				if cachedResponse.CacheTime < hourAgo {
					log.Println("Clearing cache for key:", key, "due to age")
					delete(cache, key)
				}
			}
			mu.Unlock()

			time.Sleep(1 * time.Minute) // Check every minute
		}
	}()
}
