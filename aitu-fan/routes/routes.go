package routes

import (
	"github.com/gorilla/mux"
	"github.com/group-project/aitu-fan/handlers"
	"github.com/group-project/aitu-fan/middleware"
	"net/http"
)

func SetupRoutes(router *mux.Router) {
	router.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/register.html")
	}).Methods("GET")
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/login.html")
	}).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/home.html")
	}).Methods("GET")

	router.HandleFunc("/dashboard", handlers.Dashboard).Methods("GET")
	router.HandleFunc("/create_post", handlers.CreatePostPage).Methods("GET")
	router.HandleFunc("/create_event", handlers.CreateEventPage).Methods("GET")
	router.HandleFunc("/events", handlers.EventsPage).Methods("GET")

	router.HandleFunc("/gallery", handlers.GalleryPage).Methods("GET")

	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)
	api.HandleFunc("/posts", handlers.GetPosts).Methods("GET")
	api.HandleFunc("/posts", handlers.CreatePost).Methods("POST")
	api.HandleFunc("/events", handlers.CreateEvent).Methods("POST")
	api.HandleFunc("/gallery", handlers.GetGallery).Methods("GET")
	api.HandleFunc("/gallery", handlers.UploadGalleryItem).Methods("POST")

	router.Handle("/myposts", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetMyPosts))).Methods("GET")
	router.Handle("/editpost/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.EditPostPage))).Methods("GET")
	router.Handle("/editpost/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdatePostFromForm))).Methods("POST")
	router.Handle("/deletepost/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeletePostHandler))).Methods("POST")

	admin := router.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AdminMiddleware)
	admin.HandleFunc("/users", handlers.GetUsersHandler).Methods("GET")
	admin.HandleFunc("/users/{id}", handlers.DeleteUserHandler).Methods("DELETE")

	router.Handle("/api/comments", middleware.AuthMiddleware(http.HandlerFunc(handlers.AddComment))).Methods("POST")
	router.Handle("/api/comments", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetComments))).Methods("GET")

	router.Handle("/api/like", middleware.AuthMiddleware(http.HandlerFunc(handlers.LikePost))).Methods("POST")
	router.Handle("/api/like", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetLikes))).Methods("GET")
}
