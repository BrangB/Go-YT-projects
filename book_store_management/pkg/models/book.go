package models

import (
	"log"

	"github.com/brangb/book_store_management/pkg/config"
	"gorm.io/gorm"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

func init() {
	config.Connect()    // Establish the database connection
	db = config.GetDB() // Correctly assign to the package-level variable
	db.AutoMigrate(&Book{})
}

func (b *Book) CreateBook() *Book {
	db.Create(b)
	return b
}

func GetAllBooks() []Book {
	var books []Book
	if err := db.Find(&books).Error; err != nil {
		log.Printf("Error fetching books: %v", err)
		return nil
	}
	return books
}

func GetBookById(Id int64) (*Book, *gorm.DB) {
	var getBook Book
	result := db.Where("ID=?", Id).Find(&getBook)
	if result.Error != nil {
		log.Printf("Error fetching book by ID %d: %v", Id, result.Error)
		return nil, result
	}
	return &getBook, result
}

func DeleteBook(Id int64) Book {
	var book Book
	if err := db.Where("ID=?", Id).Delete(&book).Error; err != nil {
		log.Printf("Error deleting book with ID %d: %v", Id, err)
	}
	return book
}
