package main

import (
	"log"
	"net/http"

	"go-book-api/book"
	"go-book-api/db"
	"go-book-api/login"
	"go-book-api/middleware"
	"go-book-api/post"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	
	db.InitDB()
	// Initialize mock data
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Route Handlers / Endpoints

	r.HandleFunc("/api/register", login.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", login.LoginHandler).Methods("POST")

	bookRouter := r.PathPrefix("/api/books").Subrouter()
	bookRouter.Use(middleware.JWTMiddleware)
	bookRouter.HandleFunc("", book.GetBooks).Methods("GET")
	bookRouter.HandleFunc("/{id}", book.GetBook).Methods("GET")
	bookRouter.HandleFunc("", book.CreateBook).Methods("POST")
	bookRouter.HandleFunc("/{id}", book.UpdateBook).Methods("PUT")
	bookRouter.HandleFunc("/{id}", book.DeleteBook).Methods("DELETE")

	postRouter := r.PathPrefix("/api/posts").Subrouter()
    postRouter.Use(middleware.JWTMiddleware)
    postRouter.HandleFunc("", post.CreatePost).Methods("POST")
    postRouter.HandleFunc("", post.GetPosts).Methods("GET")
    postRouter.HandleFunc("/{id}", post.DeletePost).Methods("DELETE")
    postRouter.HandleFunc("/{id}/like", post.LikePost).Methods("POST")
    postRouter.HandleFunc("/{id}/dislike", post.DislikePost).Methods("POST")
    postRouter.HandleFunc("/{id}/comment", post.CommentPost).Methods("POST")
    postRouter.HandleFunc("/{id}/favorite", post.FavoritePost).Methods("POST")

	// Start Server
	log.Println("Server starting on :8000")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS()(r)))
}
