package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// ProfileData defines the structure of profile data fetched from the backend
type ProfileData struct {
	Username     string    `json:"username"`
	Bio          string    `json:"bio"`
	ImageURL     string    `json:"image_url"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedAtSrt string    `json:"created_at_srt"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Keywords     string    `json:"keywords"`
	Author       string    `json:"author"`
	CanonicalURL string    `json:"canonical_url"`
}

// ServeProfilePage renders the profile page server-side using the fetched profile data
func ServeProfilePage(w http.ResponseWriter, r *http.Request) {
	// Extract the "username" from the URL path
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		http.Error(w, "Username not provided", http.StatusBadRequest)
		return
	}

	// Fetch profile data from the backend API
	profileData, err := fetchProfileData(username)
	if err != nil {
		http.Error(w, "Failed to fetch profile data", http.StatusInternalServerError)
		return
	}

	// Enrich profileData with template-specific fields
	profileData.Title = profileData.Username + "'s Profile"
	profileData.Description = "Learn more about " + profileData.Username + " on our site."
	profileData.Keywords = "user, profile, details"
	profileData.Author = "beauticket"
	profileData.CanonicalURL = "https://beauticket.app/" + profileData.Username
	profileData.CreatedAtSrt = profileData.CreatedAt.Format("Monday, January 2, 2006 at 3:04 PM")

	// Parse the HTML template
	tmpl, err := template.ParseFiles("template/profile.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	// Render the template with profile data
	err = tmpl.Execute(w, profileData)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// fetchProfileData fetches profile data from the backend API
func fetchProfileData(username string) (*ProfileData, error) {
	backendURL := os.Getenv("BACKEND_API_URL") + ":" + os.Getenv("PROFILE_SERVICE_PORT") + "/api/v1/profile?username=" + username

	// Make a request to the backend API
	resp, err := http.Get(backendURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var profile ProfileData
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}
