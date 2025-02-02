package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/group-project/aitu-fan/aitu-fan/handlers"
	"github.com/group-project/aitu-fan/aitu-fan/middleware"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)
	api.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Access to the protected route"))
	}).Methods("GET")

	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AdminMiddleware)
	admin.HandleFunc("/dashboard", handlers.AdminDashboardHandler).Methods("GET")
	admin.HandleFunc("/users", handlers.GetUsersHandler).Methods("GET")
	admin.HandleFunc("/users/{id}", handlers.DeleteUserHandler).Methods("DELETE")

	events := r.PathPrefix("/events").Subrouter()
	events.HandleFunc("", handlers.GetEventsHandler).Methods("GET")
	events.Handle("", middleware.AdminMiddleware(http.HandlerFunc(handlers.CreateEventHandler))).Methods("POST")

	return r
}
