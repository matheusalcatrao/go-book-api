package post

import (
	"encoding/json"
	"go-book-api/db"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Comment struct {
	PostID  uuid.UUID `json:"post_id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

type Favorite struct {
	PostID string `json:"post_id"`
	UserID string `json:"user_id"`
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// Parse the UUID from the URL parameters
	postID, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

    if err := db.Database.Model(&Post{}).Where("id = ?", postID).Update("likes", gorm.Expr("likes + ?", 1)).Error; err != nil {
        http.Error(w, "Failed to like post", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Post liked"})
}

func DislikePost(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
	// Parse the UUID from the URL parameters
	postID, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

    if err := db.Database.Model(&Post{}).Where("id = ?", postID).Update("dislikes", gorm.Expr("dislikes + ?", 1)).Error; err != nil {
        http.Error(w, "Failed to dislike post", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Post disliked"})
}

func CommentPost(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
	// Parse the UUID from the URL parameters
	postID, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
    userID := r.Context().Value("userID").(string)
    comment := r.FormValue("comment")

    newComment := Comment{
        PostID:  postID,
        UserID:  userID,
        Content: comment,
    }

    if err := db.Database.Create(&newComment).Error; err != nil {
        http.Error(w, "Failed to comment on post", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(newComment)
}

func FavoritePost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID := vars["id"]
    userID := r.Context().Value("userID").(string)

    favorite := Favorite{
        PostID: postID,
        UserID: userID,
    }

    if err := db.Database.Create(&favorite).Error; err != nil {
        http.Error(w, "Failed to favorite post", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Post favorited"})
}