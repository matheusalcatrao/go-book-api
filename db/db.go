package db

import (
	"log"
	"gorm.io/driver/postgres"
	"github.com/google/uuid"
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


func InitDB() {
	dsn := "host=localhost user=admin password=password dbname=db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Migrate the schema
	err = Database.AutoMigrate(&Book{})
	err = Database.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)

	}
}
