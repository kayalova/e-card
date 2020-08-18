package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/kayalova/e-card-catalog/helpers"
	"github.com/kayalova/e-card-catalog/models"
	"github.com/kayalova/e-card-catalog/postgres"
)

// GetAllSchools return all schools
func GetAllSchools(w http.ResponseWriter, r *http.Request) {
	schools, err := getAllSchools()

	if err != nil {
		helpers.Error("Unable to get schools", http.StatusInternalServerError, w)
		return
	}

	response, err := json.Marshal(schools)
	if err != nil {
		helpers.Error("Unable to get schools", http.StatusInternalServerError, w)
		return
	}

	w.Write(response)

}

/* ------ postgres requests ------*/
func getAllSchools() ([]models.School, error) {
	db := postgres.CreateConnection()
	defer db.Close()

	sqlStatement := `SELECT * FROM schools`
	var schools []models.School
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return schools, err
	}

	for rows.Next() {
		var school models.School
		err = rows.Scan(&school.ID, &school.Name)
		if err != nil {
			return schools, err
		}

		schools = append(schools, school)
	}

	return schools, nil
}
