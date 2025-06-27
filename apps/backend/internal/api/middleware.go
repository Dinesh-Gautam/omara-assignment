package api

import (
	"context"
	"net/http"
	"strings"

	"strategic-insight-analyst/backend/database"
	"strategic-insight-analyst/backend/firebase"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := firebase.AuthClient.VerifyIDToken(context.Background(), tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Find or create user in the database
		_, err = database.FindUserByID(token.UID)
		if err != nil {
			http.Error(w, "Failed to find user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Add user info to context
		ctx := context.WithValue(r.Context(), "user", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
