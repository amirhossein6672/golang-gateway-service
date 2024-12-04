package middleware

import (
	"context"
	"net/http"
	"golang-gateway-service/common"
)

// LoggedInMiddleware checks if the user is authenticated by validating the JWT token
func LoggedInMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract and validate admin ID from the token
		userID, err := common.ExtractUserIDFromToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Check if the user ID is valid (not zero)
		if userID <= 0 {
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		// Store the user ID in the context (if you want to use it later in the request)
		ctx := context.WithValue(r.Context(), "userID", userID)

		// Proceed to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
