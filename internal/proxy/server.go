package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func StartServer(port int, origin string) error {
	// This function will start the caching proxy server.
	// It will listen on the specified port and forward requests to the origin server.
	// The implementation details will be added later.
	// For now, we can log the port and origin for debugging purposes.
	println("Starting caching proxy server on port:", port, "with origin:", origin)

	originURL, err := url.Parse(origin)
	if err != nil {
		log.Println("Error parsing origin URL:", err)
		return fmt.Errorf("invalid origin URL: %v", err)
	}

	println("Parsed origin URL:", originURL.String())
	// Here we would set up the HTTP server, routes, and handlers.
	// This will include handling requests, checking the cache,
	// and forwarding requests to the origin server if not cached.

	proxyHandler := func(w http.ResponseWriter, r *http.Request) {
		// Check if the response is cached
		cachekey := r.Method + ":" + r.URL.String()

		if cachedResponse, found := Get(cachekey); found {
			// If cached, write the cached response
			w.WriteHeader(cachedResponse.Status)
			for key, values := range cachedResponse.Headers {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.Header().Set("X-Cache", "HIT")
			w.Write(cachedResponse.Body)
			return
		}
		// If not cached, forward the request to the origin server
		resp, err := http.Get(originURL.String() + r.URL.Path)
		if err != nil {
			http.Error(w, "Error fetching from origin", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Check for errors in the response
		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Error from origin server: "+resp.Status, resp.StatusCode)
			return
		}

		// Read the response body
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}

		// Cache the response
		Set(cachekey, &Cache{
			Body:    bodyBytes,
			Headers: resp.Header,
			Status:  resp.StatusCode,
		})

		// Write the response headers and body
		w.WriteHeader(resp.StatusCode)
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.Write(bodyBytes)
		w.Header().Set("X-Cache", "MISS")

	}

	log.Println("Setting up HTTP server on port:", port)

	return http.ListenAndServe(":"+string(port), http.HandlerFunc(proxyHandler))
}
