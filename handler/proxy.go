package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// NewBackendProxy creates a reverse proxy that routes requests based on URL path
func BackendProxy() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Determine which service to route based on the URL path
		var target *url.URL
		var err error

		// Start the HTTP server
		baseUrl := os.Getenv("BACKEND_API_URL")
		if baseUrl == "" {
			// Default to localhost
			baseUrl = "http://localhost"
		}

		switch {
		case r.URL.Path == "/api/" || startsWith(r.URL.Path, "/api/"):
			target, err = url.Parse(baseUrl + ":" + os.Getenv("BACKEND_SERVICE_PORT"))
		default:
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}

		if err != nil {
			log.Printf("⚠️ Error parsing URL: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Create a reverse proxy for the target service
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			http.Error(w, "Proxy error: "+e.Error(), http.StatusBadGateway)
		}

		// Serve the proxy
		proxy.ServeHTTP(w, r)
	})
}

// startsWith is a helper function to check if a path starts with a given prefix
func startsWith(path, prefix string) bool {
	return len(path) >= len(prefix) && path[:len(prefix)] == prefix
}
