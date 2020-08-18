package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kayalova/e-card-catalog/constants"

	"github.com/gorilla/mux"
	"github.com/kayalova/e-card-catalog/helpers"
	"github.com/kayalova/e-card-catalog/models"
	"github.com/kayalova/e-card-catalog/postgres"
)

// CreateCard creates a student card - WORKS
func CreateCard(w http.ResponseWriter, r *http.Request) {
	var card models.Card

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		helpers.Error("Unable to create the card", http.StatusConflict, w)
		return
	}

	err = insertCard(&card)

	if err != nil {
		helpers.Error("Unable to create the card", http.StatusConflict, w)
		return
	}

	w.WriteHeader(http.StatusOK)

}

// FilterCards filters students' cards
func FilterCards(w http.ResponseWriter, r *http.Request) {
	cardMap := helpers.PrepareDBfilters(r.URL.Query())
	sqlStatement := helpers.FinishUpSQLStatement(constants.SQLStatements["cards"], &cardMap)
	cards, err := filterAllRecords(sqlStatement)
	if err != nil {
		helpers.Error("Unable to get cards", http.StatusConflict, w)
		return
	}

	response := helpers.RemoveCardDuplicates(cards)
	JSONresponse, err := json.Marshal(response)
	if err != nil {
		helpers.Error("Unable to get cards", http.StatusConflict, w)
		return
	}

	w.Write(JSONresponse)

}

// EditCard updates student's card - WORKS
func EditCard(w http.ResponseWriter, r *http.Request) {
	var card models.Card
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	err = json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		helpers.Error("Unable to update the card", http.StatusInternalServerError, w)
		return

	}

	err = updateCard(id, card)
	if err != nil {
		helpers.Error("Unable to update the card", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAllCards returns all the sudents' cards - WORKS
func GetAllCards(w http.ResponseWriter, r *http.Request) {
	cards, err := getCards()
	if err != nil {
		helpers.Error("Unable to get cards", http.StatusInternalServerError, w)
		return
	}

	response := map[string][]models.Card{
		"cards": cards,
	}

	// TODO: уточнить
	JSONresponse, err := json.Marshal(response)
	if err != nil {
		helpers.Error("Unable to get cards", http.StatusInternalServerError, w)
		return
	}

	w.Write(JSONresponse)
}

// GetOneCard returns a single student's card - WORKS
func GetOneCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	card, err := getCard(id)
	if err != nil {
		helpers.Error("Unable to get the card", http.StatusInternalServerError, w)
		return
	}

	books, err := getBooksAttachedToCard(id)
	if err != nil {
		helpers.Error("Unable to get the card", http.StatusInternalServerError, w)
		return
	}

	JSONResponse := map[string]interface{}{
		"card":  card,
		"books": books,
	}

	response, err := json.Marshal(JSONResponse)
	if err != nil {
		helpers.Error("Unable to get the card", http.StatusInternalServerError, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

// DeleteOneCard deletes student's card _ WORKS
func DeleteOneCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.Error("Unable to delete the card", http.StatusInternalServerError, w)
		return
	}

	rowsCount, err := deleteCard(id)

	if err != nil || rowsCount == 0 {
		helpers.Error("Unable to delete the card", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AttachToCard atteches book to card - WORKS
func AttachToCard(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	cardID, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.Error("Unable to attach book to the card", http.StatusInternalServerError, w)
		return
	}

	bookIDStr := r.URL.Query()["id"][0]
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		helpers.Error("Unable to attach book to the card", http.StatusInternalServerError, w)
		return
	}

	err = attachBook(cardID, bookID)

	if err != nil {
		helpers.Error("Unable to attach book to the card", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// DetachFromCard detaches book from card - WORKS
func DetachFromCard(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	cardID, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.Error("Unable to detach book from the card", http.StatusInternalServerError, w)
		return
	}

	bookIDStr := r.URL.Query()["id"][0]
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		helpers.Error("Unable to detach book from the card", http.StatusInternalServerError, w)
		return
	}

	err = detachBook(cardID, bookID)

	if err != nil {
		helpers.Error("Unable to detach book from the card", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/* ------------ Postgres requests ---------- */
func updateCard(id int, card models.Card) error {
	db := postgres.CreateConnection()
	defer db.Close()

	sqlStatement := `UPDATE cards 
	SET name = $2, lastname = $3, surname= $4, phone= $5, school_id = $6	
	WHERE id = $1`

	result, err := db.Exec(sqlStatement, id, card.Name, card.Lastname, card.Surname, card.Phone, card.SchoolId)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil || count < 1 {
		return fmt.Errorf("Unable to update the card")
	}

	return nil
}

func getCards() ([]models.Card, error) {
	db := postgres.CreateConnection()
	defer db.Close()

	var cards []models.Card

	sqlStatement := `SELECT * FROM cards`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return cards, err
	}

	var card models.Card
	for rows.Next() {
		err = rows.Scan(&card.ID, &card.Name, &card.Lastname, &card.Surname, &card.Phone, &card.SchoolId)
		if err != nil {
			return cards, err
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func getCard(id int) (models.Card, error) {
	db := postgres.CreateConnection()
	defer db.Close()

	var card models.Card

	sqlStatement := `SELECT * FROM cards WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&card.ID, &card.Name, &card.Lastname, &card.Surname, &card.SchoolId, &card.Phone)
	if err != nil {
		return card, err
	}

	return card, nil
}

func insertCard(card *models.Card) error {
	db := postgres.CreateConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO cards(name, lastname, surname, phone, school_id) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(sqlStatement, card.Name, card.Lastname, card.Surname, card.Phone, card.SchoolId)

	if err != nil {
		return err
	}

	return nil
}

func deleteCard(id int) (rowsCount int64, err error) {
	db := postgres.CreateConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM cards WHERE id=$1`
	rows, err := db.Exec(sqlStatement, id)
	if err != nil {
		return 0, err
	}

	count, err := rows.RowsAffected()
	if err != nil || count < 1 {
		return 0, fmt.Errorf("Unable to delete the card")
	}

	return count, nil
}

func attachBook(cardID, bookID int) error {
	db := postgres.CreateConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO cards_books(card_id, book_id) VALUES($1, $2)`
	_, err := db.Exec(sqlStatement, cardID, bookID)
	if err != nil {
		return err
	}

	return nil
}

func detachBook(cardID, bookID int) error {
	db := postgres.CreateConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM cards_books WHERE card_id=$1 AND book_id=$2`
	_, err := db.Exec(sqlStatement, cardID, bookID)
	if err != nil {
		return err
	}

	return nil
}
