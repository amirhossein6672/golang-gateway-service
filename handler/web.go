package handler

import (
	"net/http"
	"os"
	"path/filepath"
)

// ServeWebApp serves web app files (HTML, CSS, JS) from the "static/web" directory and handles 404 errors
func ServeWebApp(webAppDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create the full path to the requested file
		filePath := filepath.Join(webAppDir, r.URL.Path[len("/web/"):])

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// Serve the custom 404 page for the web app
			http.ServeFile(w, r, filepath.Join(webAppDir, "404.html"))
			return
		}

		// If the file exists, serve it
		http.StripPrefix("/web/", http.FileServer(http.Dir(webAppDir))).ServeHTTP(w, r)
	})
}
