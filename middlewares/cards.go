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

// CreateCard creates a student card
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

// EditCard updates student's card
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

// GetAllCards returns all the sudents' cards
func GetAllCards(w http.ResponseWriter, r *http.Request) {
	cards, err := getCardsDetails()
	if err != nil {
		helpers.Error("Unable to get cards", http.StatusInternalServerError, w)
		return
	}

	JSONresponse := helpers.RemoveCardDuplicates(cards)
	response, err := json.Marshal(JSONresponse)
	if err != nil {
		helpers.Error("Unable to get cards", http.StatusInternalServerError, w)
		return
	}

	w.Write(response)
}

// GetOneCard returns a single student's card
func GetOneCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	records, err := getOneCardDetails(id)
	if err != nil {
		helpers.Error("Unable to get the card", http.StatusInternalServerError, w)
		return
	}

	JSONresponse := helpers.RemoveCardDuplicates(records)
	response, err := json.Marshal(JSONresponse)
	if err != nil {
		helpers.Error("Unable to get the card", http.StatusInternalServerError, w)
		return
	}
	w.Write(response)
}

// DeleteOneCard deletes student's card
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

// AttachToCard atteches book to card
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

// DetachFromCard detaches book from card
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

func getCardsDetails() ([]models.CommonJSON, error) {
	db := postgres.CreateConnection()
	defer db.Close()

	sqlStatement := constants.SQLStatements["cards"]
	records, err := filterAllRecords(sqlStatement)
	if err != nil {
		return make([]models.CommonJSON, 0, 0), err
	}

	return records, nil
}

func getOneCardDetails(id int) ([]models.CommonJSON, error) {

	sqlStatement := constants.SQLStatements["cards"]
	sqlStatement += fmt.Sprintf(` WHERE cards.id=%v`, id)
	records, err := filterAllRecords(sqlStatement)

	if err != nil {
		return make([]models.CommonJSON, 0, 0), err
	}

	return records, nil
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
