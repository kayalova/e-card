package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/kayalova/e-card-catalog/helpers"
	"github.com/kayalova/e-card-catalog/models"
	"github.com/kayalova/e-card-catalog/postgres"
)

// SignIn ...
func SignIn(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.Error("Unable to sign in", http.StatusInternalServerError, w)
		return
	}

	isExists, err := isEmailRegistered(user.Email)
	if err != nil {
		helpers.Error("Unable to sign in", http.StatusInternalServerError, w)
		return
	}

	if !isExists {
		helpers.Error("No user exists with such email", http.StatusUnprocessableEntity, w)
		return
	}

	registeredUser, err := findUser(&user)
	if err != nil {
		helpers.Error("Unable to sign in", http.StatusInternalServerError, w)
		return
	}

	areSamePasswords := helpers.CheckPasswordHash(user.Password, registeredUser.Password)
	if !areSamePasswords {
		helpers.Error("Wrong password", http.StatusInternalServerError, w)
		return
	}

	token, err := generateToken(int(registeredUser.ID))
	if err != nil {
		helpers.Error("Unable to sign in", http.StatusInternalServerError, w)
		return
	}

	w.Write([]byte(token))
}

// SignUp ...
func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.Error("Unable to sign up a user", http.StatusInternalServerError, w)
		return
	}

	if !helpers.IsValidUser(&user) {
		helpers.Error("Invalid data", http.StatusInternalServerError, w)
		return
	}

	isExists, err := isEmailRegistered(user.Email)
	if err != nil {
		helpers.Error("Unable to sign up a user", http.StatusInternalServerError, w)
		return
	}
	if isExists {
		helpers.Error("User with this email is already registered", http.StatusUnprocessableEntity, w)
		return
	}

	// hash pass
	bytePassword := helpers.ConvertToBytes(user.Password)
	hash := helpers.HashAndSalt(bytePassword)
	id, err := registerUser(&user, hash)

	if err != nil {
		helpers.Error("Unable to sign up a user", http.StatusInternalServerError, w)
		return
	}

	token, err := generateToken(id)
	if err != nil {
		helpers.Error("Unable to sign in", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(token))

}

// IsAuthorized checks whether user has access
func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			helpers.Error("Not authorize", http.StatusUnauthorized, w)
			return
		}

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Got an error")
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil {
			helpers.Error("Unable to handle request. Relogin please and try againg", http.StatusInternalServerError, w)
			return
		}

		if token.Valid {
			next.ServeHTTP(w, r)
		}

	})
}

/* ------ postgres requests ------ */
func generateToken(id int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = id
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	key := os.Getenv("SECRET_KEY")
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func findUser(userToFind *models.User) (models.User, error) {
	db := postgres.CreateConnection()
	defer db.Close()

	var user models.User
	sqlStatement := `SELECT * FROM users WHERE email=$1`
	row := db.QueryRow(sqlStatement, userToFind.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func registerUser(user *models.User, password string) (int, error) {
	db := postgres.CreateConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO users(name, email, password) VALUES($1, $2, $3) RETURNING id`
	id := new(int)
	err := db.QueryRow(sqlStatement, user.Name, user.Email, password).Scan(id)
	if err != nil {
		return 0, err
	}

	return *id, nil
}

func isEmailRegistered(email string) (bool, error) {
	db := postgres.CreateConnection()
	defer db.Close()

	var isExists bool
	sqlStatement := `select exists(select 1 from users where email=$1)`
	row := db.QueryRow(sqlStatement, email)
	err := row.Scan(&isExists)
	if err != nil {
		return false, err
	}

	return isExists, nil
}
