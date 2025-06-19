package routes

import (
	"net/http"

	"real-time-forum/backend/handlers"
)

func Routers() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	//  Serve static files
	fileServer := http.FileServer(http.Dir("frontend/css"))
	mux.Handle("/css/", http.StripPrefix("/css/", fileServer))
	scriptServer := http.FileServer(http.Dir("frontend/js/"))
	mux.Handle("/js/", http.StripPrefix("/js/", scriptServer))
	imageServer := http.FileServer(http.Dir("frontend/image"))
	mux.Handle("/static/image/", http.StripPrefix("/static/image/", imageServer))
	uploadServer := http.FileServer(http.Dir("./uploads/"))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", uploadServer))

	mux.HandleFunc("/register", handlers.RegisterUser)
	mux.HandleFunc("/login", handlers.HandleLogin)
	mux.HandleFunc("/logout", handlers.LogoutUser)
	
	// These endpoints require authentication
	mux.HandleFunc("/", handlers.HandleHomepage)
	mux.HandleFunc("/users", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandleUsers)))
	mux.HandleFunc("/upload-image", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandleImageUpload)))
	mux.HandleFunc("/categories", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandleGetCategories)))
	mux.HandleFunc("/comments", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandleComments)))
	mux.HandleFunc("/posts", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandlePosts)))
	mux.HandleFunc("/ws", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandleWebSocket)))
	mux.HandleFunc("/messages", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandleGetMessages)))
	mux.HandleFunc("/posts/like", handlers.AuthMiddleware(http.HandlerFunc(handlers.LikePostHandler)))
	mux.HandleFunc("/posts/dislike", handlers.AuthMiddleware(http.HandlerFunc(handlers.DislikePostHandler)))
	mux.HandleFunc("/comments/like", handlers.AuthMiddleware(http.HandlerFunc(handlers.LikeCommentHandler)))
	mux.HandleFunc("/comments/dislike", handlers.AuthMiddleware(http.HandlerFunc(handlers.DislikeCommentHandler)))

	return mux, nil
}
