package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func GoogleMapProxy() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract the relevant portion of the path after "googleapis/"
        pathAfterGoogleapis := strings.TrimPrefix(r.URL.Path, "/googleapis/")
        
        // Base target URL
        targetBase := "https://maps.googleapis.com/"
        
        // Construct the target URL
        targetURL, err := url.Parse(targetBase)
        if err != nil {
            log.Printf("⚠️ Error parsing target base URL: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        
        // Update target URL's Path and Query
        targetURL.Path += pathAfterGoogleapis
        targetURL.RawQuery = r.URL.RawQuery + "&key=" + os.Getenv("GOOGLE_MAP_WEB_APIKEY")

        log.Printf("Proxying request to: %s", targetURL.String())
        
        // Create and serve the reverse proxy
        proxy := httputil.NewSingleHostReverseProxy(targetURL)
        proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
            log.Printf("Proxy error: %v", err)
            http.Error(w, "Proxy Error: "+err.Error(), http.StatusBadGateway)
        }

        proxy.ServeHTTP(w, r)
    })
}
