package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kayalova/e-card-catalog/helpers"
	"github.com/kayalova/e-card-catalog/models"
	"github.com/kayalova/e-card-catalog/postgres"
	_ "github.com/lib/pq"
)

// Create a student card
func Create(w http.ResponseWriter, r *http.Request) {
	var card models.Card

	err := json.NewDecoder(r.Body).Decode(&card)
	// роняет сервер
	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
		return
	}

	err = insertCard(&card)

	if err != nil {
		// роняет сервер
		log.Fatalf("Unable to execute the query. %v", err)
		// check later
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Unable to create the card"))
	} else {
		w.WriteHeader(http.StatusOK)
	}

}

// Filter students' cards
func FilterCards(w http.ResponseWriter, r *http.Request) {

	filters := r.URL.Query()
	cardMap := helpers.PrepareDBfilters(filters)

	var card models.Card
	sqlStatement := helpers.BuildSQLStatement(&cardMap, "cards", models.Card)

	cards, error := helpers.GetDBRecords(sqlStatement, card)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to get the cards"))
		return
	}

	var cardsID []int
	for _, card := range cards {
		cardsID = append(cardsID, card.ID)
	}

	booksSqlSatement = fmt.Sprintf("SELECT book_id FROM cards_books WHERE card_id IN(%v)", cardsID)
	// TODO: придумать адекватное решение
	strings.Replace(booksSqlSatement, "[", "(")
	strings.Replace(booksSqlSatement, "]", ")")

	var book models.Book
	books, error := helpers.GetDBRecords(booksSqlSatement, book, models.Book)

	var JSONResponse = map[string]interface{}{
		"cards": cards,
		"books": books,
	}

	response, err := json.Marshal(JSONResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to get cards"))
		return
	}

	w.Write(response)
	w.Header().Set("Content-Type", "application/json")

}

// EditCard updates student's card
func EditCard(w http.ResponseWriter, r *http.Request) {
	var card models.Card
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	err = json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	err = updateCard(id, card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to update the card"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAllCards returns all the sudents' cards
func GetAllCards(w http.ResponseWriter, r *http.Request) {
	cards, err := getCards()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to get cards"))
		return
	}

	JSONResponse := map[string][]models.Card{
		"cards": cards,
	}

	response, err := json.Marshal(JSONResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to get cards"))
		return
	}

	w.Write(response)
	w.Header().Set("Content-Type", "application/json")
}

// GetOneCard returns a single student's card
func GetOneCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	card, err := getCard(id)
	// TODO: пофиксить эту дичь с ошибками
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to get the card"))
		return
	}

	books, err := getBooksAttachedToCard(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to get the card"))
		return
	}

	JSONResponse := map[string]interface{}{
		"card":  card,
		"books": books,
	}

	response, err := json.Marshal(JSONResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to get the card"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

// DeleteOneCard deletes student's card
func DeleteOneCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// роняет сервер
		log.Fatalf("Got an error while parsing: %v", err)
	}

	rowsCount, err := deleteCard(id)

	if err != nil || rowsCount == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to delete the card"))
	} else {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Successfully deleted"))
	}
}

/* ------------ Postgres requests ---------- */
func findAll(filterMap map[string]interface{}, tableName string) (interface{}, error) {

	return records, nil
}

func updateCard(id int, card models.Card) error {
	db := postgres.CreateConnection()
	defer db.Close()

	fmt.Println(card)

	sqlStatement := `UPDATE cards 
	SET name = $2, lastname = $3, surname= $4, phone= $5, school_id = $6	
	WHERE id = $1`

	_, err := db.Exec(sqlStatement, id, card.Name, card.Lastname, card.Surname, card.Phone, card.SchoolId)
	if err != nil {
		return err
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
	if err != nil {
		return 0, err
	}

	return count, nil
}
