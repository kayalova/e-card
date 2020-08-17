package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/kayalova/e-card-catalog/models"
)

// PrepareDBfilters ...
func PrepareDBfilters(filters url.Values) map[string]interface{} {
	mapFilters := make(map[string]interface{})
	for k, v := range filters {
		mapFilters[k] = v[0]
	}

	return mapFilters
}

// BuildSQLStatement returns sql statement string - NOT USED
func BuildSQLStatement(m *map[string]interface{}, tableName string) string {
	sqlStatement := fmt.Sprintf(`SELECT * FROM %v WHERE `, tableName)
	for k, v := range *m {
		sqlStatement += fmt.Sprintf(`%v=%v AND `, k, v)
	}

	return strings.ReplaceAll(sqlStatement[:len(sqlStatement)-4], `"`, `'`)
}

// Error ...
func Error(msg string, httpCode int, w http.ResponseWriter) {
	var response models.Response = models.Response{
		Message: msg,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		Error("Что-то пошло не так", httpCode, w)
	} else {
		w.WriteHeader(httpCode)
		w.Write(responseJSON)
	}
}

// CreateCardsSQLStatement ...
func CreateCardsSQLStatement(filters *map[string]interface{}) string {
	sqlStatement := `SELECT
	cards.id card_id,
	cards.name card_name,
	cards.lastname card_lastname,
	cards.surname card_surname,
	cards.phone card_phone,
	schools.name school_name,
	books.id book_id,
	books.name book_name,
	books.author book_author,
	books.book_id book_ownID
  FROM
	cards
  LEFT JOIN schools
	ON schools.id = cards.school_id
  LEFT JOIN cards_books
	ON cards_books.card_id = cards.id
  LEFT JOIN books
	ON books.id = cards_books.book_id;
	`

	if len(*filters) == 0 {
		return sqlStatement
	}

	sqlStatement += `WHERE `
	for k, v := range *filters {
		sqlStatement += fmt.Sprintf(`%v=%v AND `, k, v)
	}

	return strings.ReplaceAll(sqlStatement[:len(sqlStatement)-4], `"`, `'`)

}

// RemoveDuplicates ...
func RemoveDuplicates(records []models.CommonJSON) map[int64]map[string]interface{} {

	setID := make(map[int64]bool)
	result := make(map[int64]map[string]interface{})

	for _, m := range records {
		if setID[m.Card.ID] {
			books := result[m.Card.ID]["Books"]
			books = append(books.([]interface{}), m.Book)
			result[m.Card.ID]["Books"] = books
		} else {
			result[m.Card.ID] = map[string]interface{}{
				"Card":   m.Card,
				"School": m.School,
				"Books":  make([]interface{}, 0, 5),
			}

			if !IsEmptyString(m.Book.Name) {
				books := result[m.Card.ID]["Books"]
				books = append(books.([]interface{}), m.Book)
				result[m.Card.ID]["Books"] = books
			}

			setID[m.Card.ID] = true
		}
	}

	return result
}

// IsEmptyString ...
func IsEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}
