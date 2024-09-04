package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"go-book-api/db"

)

// Define your JWT secret key
var jwtKey = []byte("my_secret_key")

// User struct represents a user in the database.
type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Username string    `gorm:"unique;not null" json:"username"`
	Password string    `gorm:"not null" json:"password"`
}

// UserCredentials represents the data required for login.
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims struct for JWT token claims.
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// RegisterHandler handles user registration.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Hash the password before storing it.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Create user in the database
	if result := db.Database.Create(&user); result.Error != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User registered successfully")
}


// LoginHandler handles user login.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials UserCredentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user User
	// Check if user exists in the database
	if result := db.Database.Where("username = ?", credentials.Username).First(&user); result.Error != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare stored hashed password with the provided one
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create a JWT token
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	// Send token as JSON
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
