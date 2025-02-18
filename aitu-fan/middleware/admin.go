package middleware

import (
	"fmt"
	"net/http"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAdmin(r) {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		fmt.Println("Admin access granted")
		next.ServeHTTP(w, r)
	})
}

func GetUserRole(r *http.Request) (string, error) {
	userRole, ok := r.Context().Value("userRole").(string)
	if !ok {
		return "", fmt.Errorf("Error: unable to get user role")
	}
	return userRole, nil
}

func IsAdmin(r *http.Request) bool {
	userRole, err := GetUserRole(r)
	if err != nil {
		return false
	}
	return userRole == "admin"
}
