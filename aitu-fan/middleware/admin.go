package middleware

import (
	"fmt"
	"github.com/group-project/aitu-fan/aitu-fan/utils"
	"net/http"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Authorization token missing", http.StatusUnauthorized)
			return
		}

		user, err := utils.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if user.Role != "admin" {
			http.Error(w, "Insufficient permissions", http.StatusForbidden)
			return
		}

		fmt.Println("Admin access granted to:", user.Username)
		next.ServeHTTP(w, r)
	})
}
