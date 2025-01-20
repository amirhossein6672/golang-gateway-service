package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// NewFileProxy creates a reverse proxy that routes requests based on URL path
func FileProxy() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { 
		// Determine which service to route based on the URL path
		var target *url.URL
		var err error

		// Start the HTTP server
		baseUrl := os.Getenv("FILE_API_URL")
		if baseUrl == "" {
			// Default to localhost
			baseUrl = "http://localhost"
		}

		switch {
		case r.URL.Path == "/" || startsWith(r.URL.Path, "/"):
			target, err = url.Parse(baseUrl + ":" + os.Getenv("FILE_SERVICE_PORT"))
		default:
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}

		if err != nil {
			log.Printf("⚠️ Error parsing URL: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Printf("Sending request to: %s", target.String())

		// Create a reverse proxy for the target service
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			http.Error(w, "Proxy error: "+e.Error(), http.StatusBadGateway)
		}

		// Serve the proxy
		proxy.ServeHTTP(w, r)
	})
} 