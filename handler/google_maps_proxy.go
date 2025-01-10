package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// GoogleMapProxy creates a reverse proxy that routes requests based on URL path
func GoogleMapProxy() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathAfterGoogleapis := strings.TrimPrefix(r.URL.Path, "/googleapis/")
		
		// Construct the target URL
		targetBase := "https://maps.googleapis.com/"
		targetURL, err := url.Parse(targetBase + pathAfterGoogleapis)
		if err != nil {
			log.Printf("⚠️ Error parsing URL: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Copy query parameters from the original request
		query := r.URL.Query()
		query.Set("key", os.Getenv("GOOGLE_MAP_WEB_APIKEY")) // Add the API key
		targetURL.RawQuery = query.Encode()

		// Create a reverse proxy for the target service
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			log.Printf("⚠️ Proxy error: %v", e)
			http.Error(w, "Proxy error: "+e.Error(), http.StatusBadGateway)
		}

		// Update the request's URL and Host to point to the target
		r.URL = targetURL
		r.Host = targetURL.Host

		// Serve the proxy
		proxy.ServeHTTP(w, r)
	})
}
