package middlewares

import (
	"net/http"

	"github.com/kayalova/e-card-catalog/models"
	"github.com/kayalova/e-card-catalog/postgres"
)

// GetAllBooks ...
func GetAllBooks(w http.ResponseWriter, r *http.Request) {}

// FilterBooks ...
func FilterBooks(w http.ResponseWriter, r *http.Request) {}

/* ------------ Postgres requests ---------- */
func getBooksAttachedToCard(id int) ([]models.Book, error) {
	db := postgres.CreateConnection()
	defer db.Close()

	var books []models.Book
	sqlStatement := `SELECT book_id FROM cards_books WHERE card_id=$1`
	rows, err := db.Query(sqlStatement, id) // все записи с кард_ид = ид
	if err != nil {
		return books, err
	}

	var booksID []int
	var bookID int

	for rows.Next() {
		err = rows.Scan(&bookID)
		if err != nil {
			return books, err
		}
		booksID = append(booksID, bookID)
	}

	for _, bookID := range booksID {
		book, err := getBook(bookID)

		if err != nil {
			return books, err
		}

		books = append(books, book)
	}

	return books, nil
}

func getBook(id int) (models.Book, error) {
	db := postgres.CreateConnection()
	defer db.Close()

	var book models.Book
	sqlStatement := `SELECT * FROM books WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&book.ID, &book.Name, &book.Author, &book.BookId)
	if err != nil {
		return book, err
	}
	return book, nil
}
