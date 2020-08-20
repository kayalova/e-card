package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/kayalova/e-card-catalog/models"
	"golang.org/x/crypto/bcrypt"
)

// Error handler for all errors' cases
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

// CheckPasswordHash ...
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// IsValidUser ...
func IsValidUser(user *models.User) bool {
	return !IsEmptyString(user.Email) &&
		!IsEmptyString(user.Name) &&
		!IsEmptyString(user.Password)
}

// HashAndSalt ...
func HashAndSalt(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// PrepareDBfilters returns map[query_paramName]query_paramValue
func PrepareDBfilters(filters url.Values) map[string]interface{} {
	mapFilters := make(map[string]interface{})
	for k, v := range filters {
		mapFilters[k] = v[0]
	}

	return mapFilters
}

// FinishUpSQLStatement builds sql query string depends on PrepareDBfilters map
func FinishUpSQLStatement(sqlStatement string, filters *map[string]interface{}) string {
	if len(*filters) == 0 {
		return sqlStatement
	}

	sqlStatement += ` WHERE `
	for k, v := range *filters {
		sqlStatement += fmt.Sprintf(`%v=%v AND `, k, v)
	}

	return strings.ReplaceAll(sqlStatement[:len(sqlStatement)-4], `"`, `'`)
}

// RemoveCardDuplicates removes all cards' duplicates
func RemoveCardDuplicates(records []models.CommonJSON) map[int64]map[string]interface{} {
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

// ConvertToBytes converts string into []byte
func ConvertToBytes(s string) []byte {
	return []byte(s)
}
