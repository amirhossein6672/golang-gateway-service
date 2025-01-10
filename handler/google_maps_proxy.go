package handler

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
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

		log.Printf("Sending request to: %s", targetURL.String())

		// Create a new HTTP request to the target API
		req, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
		if err != nil {
			log.Printf("⚠️ Error creating request: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Copy headers from the original request
		for key, values := range r.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Perform the request using an HTTP client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("⚠️ Error performing request: %v", err)
			http.Error(w, "Failed to fetch data from Google API", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copy the response status code
		w.WriteHeader(resp.StatusCode)

		// Copy the response headers
		for key, values := range resp.Header {
			// Skip Content-Encoding header if gzip to avoid double compression
			if key == "Content-Encoding" && values[0] == "gzip" {
				continue
			}
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// Handle gzip-encoded responses
		var reader io.ReadCloser
		if resp.Header.Get("Content-Encoding") == "gzip" {
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				log.Printf("⚠️ Error decompressing response: %v", err)
				http.Error(w, "Failed to decompress response", http.StatusInternalServerError)
				return
			}
			defer reader.Close()
		} else {
			reader = resp.Body
		}

		// Copy the response body
		_, err = io.Copy(w, reader)
		if err != nil {
			log.Printf("⚠️ Error copying response body: %v", err)
		}
	})
}
