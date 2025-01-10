package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// NewBackendProxy creates a reverse proxy that routes requests based on URL path
func GoogleMapProxy() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Determine which service to route based on the URL path
		var target *url.URL
		var err error

		// Start the HTTP serve
		fullPathString := strings.Split(r.URL.Path, "googleapis/")[1]
		log.Printf(

		target, err = url.Parse("https://maps.googleapis.com/" + fullPathString + "&key=" + os.Getenv("GOOGLE_MAP_WEB_APIKEY"))

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
