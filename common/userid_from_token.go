package common

import (
	"fmt"
	"net/http"
	"golang-gateway-service/model"
	"strconv"
	"strings"
)

// ExtractUserIDFromToken extracts and validates the user ID from the Authorization header
func ExtractUserIDFromToken(r *http.Request) (int, error) {
	// Extract the token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("⚠️ %s: Authorization header missing", string(model.ErrInvalidAccessToken))
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
	if tokenString == "" {
		return 0, fmt.Errorf("⚠️ %s: Bearer token missing", string(model.ErrInvalidAccessToken))
	}

	// Validate the token and get the claims
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return 0, fmt.Errorf("⚠️ %s: Invalid token: %v", string(model.ErrInvalidAccessToken), err)
	}

	// Extract user ID from the claims
	userIDStr := claims.UserID
	if userIDStr == "0" {
		return 0, fmt.Errorf("⚠️ %s: Invalid user ID", string(model.ErrInvalidUserId))
	}

	// Convert userIDStr from string to int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, fmt.Errorf("⚠️ %s: Invalid user ID format", string(model.ErrInvalidUserId))
	}

	return userID, nil
}
