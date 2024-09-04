package main

import (
	"log"
	"net/http"

	"go-book-api/book"
	"go-book-api/login"
	"go-book-api/db"
	"go-book-api/middleware"

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


	// Start Server
	log.Println("Server starting on :8000")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS()(r)))
}
