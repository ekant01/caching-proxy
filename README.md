# Caching Proxy

A CLI-based caching HTTP proxy server written in Go. It forwards incoming HTTP requests to an origin server and caches the responses. If the same request is made again, the proxy returns the cached response instead of forwarding the request again, reducing latency and load on the origin server.
