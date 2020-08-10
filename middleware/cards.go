package middleware

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/kayalova/e-card-catalog/models"
	_ "github.com/lib/pq"
)

func createConnection() *sql.DB {
	connStr := "user=postgres password=1 dbname=e-catalog sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	return db

}

// Create a student card
func Create(w http.ResponseWriter, r *http.Request) {
	var card models.Card

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	err = insertCard(&card)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
		// check later
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Unable to create a card"))
	} else {
		w.WriteHeader(http.StatusOK)
	}

}

func Filter(w http.ResponseWriter, r *http.Request) {

}

func Edit(w http.ResponseWriter, r *http.Request) {

}

func GetAll(w http.ResponseWriter, r *http.Request) {

}

func GetOne(w http.ResponseWriter, r *http.Request) {

}

func DeleteOne(w http.ResponseWriter, r *http.Request) {

}

func insertCard(card *models.Card) error {
	db := createConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO cards(name, lastname, surname, phone, school_id) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(sqlStatement, card.Name, card.Lastname, card.Surname, card.Phone, card.SchoolId)

	if err != nil {
		return err
	}

	return nil
}
