package db

import (
	"log"
	"os"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

// Book struct represents a book with an ID, Title, Author, and Year.
type Book struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	Title  string    `json:"title"`
	Author string    `json:"author"`
	Year   string    `json:"year"`
	Photo  string    `json:"photo"`
}

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key" json:"id"`
	Username string    `gorm:"unique;not null" json:"username"`
	Password string    `gorm:"not null" json:"password"`
}

type Post struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
    UserID string `json:"user_id"`
    BookID string `json:"book_id"`
    Title  string `json:"title"`
    Body   string `json:"body"`
	Likes   int `json:"like"`
}

type Comment struct {
	PostID  string `json:"post_id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

type Favorite struct {
	PostID string `json:"post_id"`
	UserID string `json:"user_id"`
}


func InitDB() {
	// dsn := "host=localhost user=admin password=password dbname=db port=5432 sslmode=disable TimeZone=Brazil/East"
	dsn := os.Getenv("DATABASE_URL")
	var err error
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Migrate the schema
	err = Database.AutoMigrate(&Book{})
	err = Database.AutoMigrate(&User{})
	err = Database.AutoMigrate(&Post{})
	err = Database.AutoMigrate(&Comment{})
	err = Database.AutoMigrate(&Favorite{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)

	}
}
