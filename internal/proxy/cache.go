package proxy

import (
	"log"
	"net/http"
	"sync"
)

type Cache struct {
	Body    []byte
	Headers http.Header
	Status  int
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
