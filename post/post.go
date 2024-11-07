package post

import (
	"encoding/json"
	"go-book-api/db"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Post struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
    UserID string `json:"user_id"`
    BookID string `json:"book_id"`
    Title  string `json:"title"`
    Body   string `json:"body"`
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
    var post Post
	print("Creating post")
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	post.ID = uuid.New()
	
    // Save post to the database
	if result := db.Database.Create(&post); result.Error != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// Parse the UUID from the URL parameters
	postID, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	// Delete the post from the database
	if result := db.Database.Where("id = ?", postID).Delete(&Post{}); result.Error != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted"})
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
    var posts []Post	
	if result := db.Database.Find(&posts); result.Error != nil {
		http.Error(w, "Failed to get post", http.StatusInternalServerError)
		return
	}
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(posts)
}