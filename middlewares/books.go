package middlewares

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kayalova/e-card-catalog/constants"
	"github.com/kayalova/e-card-catalog/helpers"
	"github.com/kayalova/e-card-catalog/models"
	"github.com/kayalova/e-card-catalog/postgres"
)

// GetAllBooks ...
func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := getBooks()
	if err != nil {
		helpers.Error("Unable to get books", http.StatusInternalServerError, w)
		return
	}

	response, err := json.Marshal(books)
	if err != nil {
		helpers.Error("Unable to get books", http.StatusInternalServerError, w)
		return
	}

	w.Write(response)
}

// FilterBooks ...
func FilterBooks(w http.ResponseWriter, r *http.Request) {
	booksFilters := helpers.PrepareDBfilters(r.URL.Query())
	sqlStatement := helpers.FinishUpSQLStatement(constants.SQLStatements["books"], &booksFilters)
	booksAndCards, err := filterAllRecords(sqlStatement)
	if err != nil {
		helpers.Error("Unable to get books1", http.StatusInternalServerError, w)
		return
	}

	response := helpers.RemoveCardDuplicates(booksAndCards)
	JSONresponse, err := json.Marshal(response)
	if err != nil {
		helpers.Error("Unable to get cards", http.StatusConflict, w)
		return
	}

	w.Write(JSONresponse)

}

/* ------------ Postgres requests ---------- */
func getBooks() ([]models.Book, error) {
	db := postgres.CreateConnection()
	defer db.Close()

	var books []models.Book
	sqlStatement := `SELECT * FROM books`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return books, err
	}

	for rows.Next() {
		var book models.Book
		err = rows.Scan(&book.ID, &book.Name, &book.Author, &book.BookId)
		if err != nil {
			return books, err
		}
		books = append(books, book)
	}

	return books, nil
}

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

func filterAllRecords(sqlStatement string) ([]models.CommonJSON, error) {
	db := postgres.CreateConnection()
	defer db.Close()

	rows, err := db.Query(sqlStatement)

	var records []models.CommonJSON
	if err != nil {
		log.Println(err)
		return records, err
	}

	for rows.Next() {

		var complete models.CommonJSON
		err = rows.Scan(&complete.Card.ID, &complete.Card.Name, &complete.Card.Lastname, &complete.Card.Surname, &complete.Card.Phone, &complete.School.ID, &complete.School.Name, &complete.Book.ID, &complete.Book.Name, &complete.Book.Author, &complete.Book.BookId)
		if err != nil {
			log.Println(err)
		}
		records = append(records, complete)
	}

	return records, nil
}
