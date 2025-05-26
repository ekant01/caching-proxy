package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

var originURL *url.URL

func StartServer(port int, origin string) error {
	log.Println("Starting caching proxy server on port:", port, "with origin:", origin)

	var err error
	originURL, err = url.Parse(origin)
	if err != nil {
		log.Println("Error parsing origin URL:", err)
		return fmt.Errorf("invalid origin URL: %v", err)
	}

	log.Println("Setting up HTTP server on port:", port)
	return http.ListenAndServe(":"+fmt.Sprint(port), http.HandlerFunc(proxyHandler))
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request:", r.Method, r.URL.String())
	cachekey := r.Method + ":" + r.URL.String()

	if r.URL.Path == "/clear-cache" {
		log.Println("Received request to clear cache")
		ClearAll()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cache cleared successfully"))
		return
	}

	if cachedResponse, found := Get(cachekey); found {
		log.Println("Cache hit for key:", cachekey)
		for key, values := range cachedResponse.Headers {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.Header().Set("X-Cache", "HIT")
		w.WriteHeader(cachedResponse.Status)
		w.Write(cachedResponse.Body)
		return
	}

	fullURL := originURL.ResolveReference(r.URL).String()
	resp, err := http.Get(fullURL)
	if err != nil {
		http.Error(w, "Error fetching from origin", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error from origin server: "+resp.Status, resp.StatusCode)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	Set(cachekey, &Cache{
		Body:    bodyBytes,
		Headers: resp.Header,
		Status:  resp.StatusCode,
	})

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.Header().Set("X-Cache", "MISS")
	w.WriteHeader(resp.StatusCode)
	w.Write(bodyBytes)
}
