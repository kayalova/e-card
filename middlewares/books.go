package middlewares

import (
	"github.com/kayalova/e-card-catalog/models"
)

//переместить
func getBooksAttachedToCard(id int) ([]models.Book, error) {
	db := CreateConnection()
	defer db.Close()

	var books []models.Book
	sqlStatement := `SELECT book_id FROM cards_books WHERE card_id=$1`
	rows, err := db.Query(sqlStatement, id) // все записи с кард_ид = ид
	var booksID []int
	if err != nil {
		return books, err
	}

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

//переместить
func getBook(id int) (models.Book, error) {
	db := CreateConnection()
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

func GetAllBooks() {}

func FilterBooks() {}
