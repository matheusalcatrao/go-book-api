package book

import (
	"encoding/json"
	"fmt"
	"go-book-api/db"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Book struct represents a book with an ID, Title, Author, and Year.
type Book struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Title  string    `json:"title"`
	Author string    `json:"author"`
	Year   string    `json:"year"`
	Photo  string    `json:"photo"`
}

// Books slice to seed book data.
var books []Book

// GetBooks responds with the list of all books as JSON.
func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if result := db.Database.Find(&books); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(books)
}


// GetBook responds with a single book by ID.
func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get the parameters from the request URL
	id, _ := uuid.Parse(params["id"])
	for _, item := range books {
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// CreateBook adds a new book to the list.
func CreateBook(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	errEnv := godotenv.Load()
	if errEnv != nil {
	  log.Fatal("Error loading .env file")
	}

	// Parse the multipart form data to handle file uploads
	err := request.ParseMultipartForm(10 << 20) // Max upload size set to 10MB
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	// Extract book data from the form
	var book Book
	book.Title = request.FormValue("title")
	book.Author = request.FormValue("author")
	book.Year = request.FormValue("year")

	// Handle file upload
	file, handler, err := request.FormFile("photo")
	if err == nil {
		defer file.Close()

		// Create a directory to store uploaded photos if it doesn't exist
		if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
			err = os.Mkdir("./uploads", os.ModePerm)
			if err != nil {
				http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
				return
			}
		}

		baseURL := os.Getenv("BASE_URL") 
		// Generate a unique filename using UUID to avoid collisions
		fileExtension := filepath.Ext(handler.Filename)
		newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExtension)
		filePath := filepath.Join("uploads", newFileName)
		fullFileURL := fmt.Sprintf("%s/%s", baseURL, filePath)
		print(fullFileURL)
		
		// Create a new file in the uploads directory
		destFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Unable to create file", http.StatusInternalServerError)
			return
		}
		defer destFile.Close()

		// Copy the uploaded file to the new file
		if _, err := destFile.ReadFrom(file); err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		// Set the Photo field to the file path
		book.Photo = fullFileURL
	}

	// Generate a new UUID for the book ID
	book.ID = uuid.New()

	// Save the book to the database
	if result := db.Database.Create(&book); result.Error != nil {
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}

	// Respond with the created book data
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}
// UpdateBook updates an existing book.
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// Parse the UUID from the URL parameters
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	// Find the book by ID in the database
	var book Book
	if result := db.Database.First(&book, "id = ?", id); result.Error != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Decode the request body into the book struct
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the book in the database
	if result := db.Database.Save(&book); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the updated book
	json.NewEncoder(w).Encode(book)
}

// DeleteBook removes a book from the list.
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// Parse the UUID from the URL parameters
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	// Find the book by ID in the database
	var book Book
	if result := db.Database.First(&book, "id = ?", id); result.Error != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Delete the book from the database
	if result := db.Database.Delete(&book); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusNoContent)
}

func InitializeBooks() {
	books = append(books, Book{ID: uuid.New(), Title: "Golang", Author: "John Doe", Year: "2012"})
	books = append(books, Book{ID: uuid.New(), Title: "Python", Author: "Jane Doe", Year: "2013"})
}