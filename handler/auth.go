package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/kayalova/e-card-catalog/helper"
	"github.com/kayalova/e-card-catalog/model"
	"github.com/kayalova/e-card-catalog/settings"
)

// SignIn ...
func SignIn(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helper.Error("Unable to sign in", http.StatusInternalServerError, w)
		return
	}

	isExists, err := isEmailRegistered(user.Email)
	if err != nil {
		helper.Error("Unable to sign in", http.StatusInternalServerError, w)
		return
	}

	if !isExists {
		helper.Error("No user exists with such email", http.StatusUnprocessableEntity, w)
		return
	}

	registeredUser, err := findUser(&user)
	if err != nil {
		helper.Error("Unable to sign in", http.StatusInternalServerError, w)
		return
	}

	areSamePasswords := helper.CheckPasswordHash(user.Password, registeredUser.Password)
	if !areSamePasswords {
		helper.Error("Wrong password", http.StatusInternalServerError, w)
		return
	}

	token, err := generateToken(int(registeredUser.ID))
	if err != nil {
		helper.Error("Unable to sign in", http.StatusInternalServerError, w)
		return
	}

	w.Write([]byte(token))
}

// SignUp ...
func SignUp(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helper.Error("Unable to sign up a user", http.StatusInternalServerError, w)
		return
	}

	if !helper.IsValidUser(&user) {
		helper.Error("Invalid data", http.StatusInternalServerError, w)
		return
	}

	isExists, err := isEmailRegistered(user.Email)
	if err != nil {
		helper.Error("Unable to sign up a user", http.StatusInternalServerError, w)
		return
	}
	if isExists {
		helper.Error("User with this email is already registered", http.StatusUnprocessableEntity, w)
		return
	}

	// hash pass
	bytePassword := helper.ConvertToBytes(user.Password)
	hash := helper.HashAndSalt(bytePassword)
	id, err := registerUser(&user, hash)

	if err != nil {
		helper.Error("Unable to sign up a user", http.StatusInternalServerError, w)
		return
	}

	token, err := generateToken(id)
	if err != nil {
		helper.Error("Unable to sign in", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(token))

}

/* ------ postgres requests ------ */
func generateToken(id int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = id
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	key := settings.GetEnvKey("SECRET_KEY", "MY_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func findUser(userToFind *model.User) (model.User, error) {
	db := settings.CreateConnection()
	defer db.Close()

	var user model.User
	sqlStatement := `SELECT * FROM users WHERE email=$1`
	row := db.QueryRow(sqlStatement, userToFind.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func registerUser(user *model.User, password string) (int, error) {
	db := settings.CreateConnection()
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
	db := settings.CreateConnection()
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
