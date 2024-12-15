package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"golang-gateway-service/common"
	"golang-gateway-service/common/middleware"
	"golang-gateway-service/handler"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	// Ensure the "logs" directory exists; create it if it doesn't
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create log directory: %v", err)
		}
	}

	// Create a new log file with a timestamp in its name
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := filepath.Join(logDir, fmt.Sprintf("gateway_service_%s.log", timestamp))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Set up multi-writer for logging to both terminal and log file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	// Print ASCII art banner to the log
	log.Println(`
     ██████╗  █████╗ ████████╗███████╗██╗    ██╗ █████╗ ██╗   ██╗
    ██╔════╝ ██╔══██╗╚══██╔══╝██╔════╝██║    ██║██╔══██╗╚██╗ ██╔╝
    ██║  ███╗███████║   ██║   █████╗  ██║ █╗ ██║███████║ ╚████╔╝ 
    ██║   ██║██╔══██║   ██║   ██╔══╝  ██║███╗██║██╔══██║  ╚██╔╝  
    ╚██████╔╝██║  ██║   ██║   ███████╗╚███╔███╔╝██║  ██║   ██║   
     ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚══════╝ ╚══╝╚══╝ ╚═╝  ╚═╝   ╚═╝    
	`)

	

	// Load environment variables
	common.LoadEnv()

	// Get the server port from environment variables or use the default
	port := os.Getenv("GATEWAY_SERVICE_PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize the router
	router := mux.NewRouter()

	// Apply the logging middleware
	router.Use(middleware.Logging) 

	// API reverse proxy
	router.PathPrefix("/").Handler(handler.BackendProxy())

	// Server-Side Rendering for profile page, capture the username in the URL
	router.HandleFunc("/@{username}", handler.ServeProfilePage)

	// Custom NotFoundHandler to serve 404.html page
	router.NotFoundHandler = http.HandlerFunc(handler.NotFoundHandler)

	// Set up CORS middleware with desired options
	c := cors.New(cors.Options{
		// AllowedOrigins defines which origins are permitted to access the resources.
		// "*" allows all origins; in production, specify the exact origins.
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"*"},
		// AllowCredentials indicates whether the request can include user credentials like cookies.
		AllowCredentials: true,
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(router)

	// Start the server
	log.Println("🚀 Gateway service is running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
