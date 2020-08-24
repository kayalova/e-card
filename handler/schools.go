package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kayalova/e-card-catalog/helper"
	"github.com/kayalova/e-card-catalog/model"
	"github.com/kayalova/e-card-catalog/settings"
)

// GetAllSchools returns all schools
func GetAllSchools(w http.ResponseWriter, r *http.Request) {
	schools, err := getAllSchools()

	if err != nil {
		helper.Error("Unable to get schools", http.StatusInternalServerError, w)
		return
	}

	response, err := json.Marshal(schools)
	if err != nil {
		helper.Error("Unable to get schools", http.StatusInternalServerError, w)
		return
	}

	w.Write(response)

}

/* ------ postgres requests ------*/
func getAllSchools() ([]model.School, error) {
	db := settings.CreateConnection()
	defer db.Close()

	sqlStatement := `SELECT * FROM schools`
	var schools []model.School
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return schools, err
	}

	for rows.Next() {
		var school model.School
		err = rows.Scan(&school.ID, &school.Name)
		if err != nil {
			return schools, err
		}

		schools = append(schools, school)
	}

	return schools, nil
}
