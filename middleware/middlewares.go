package middleware

import (
	"fmt"
	"net/http"

	"github.com/kayalova/e-card-catalog/settings"

	"github.com/dgrijalva/jwt-go"
	"github.com/kayalova/e-card-catalog/helper"
)

// IsAuthorized checks whether user has access
func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			helper.Error("Not authorize", http.StatusUnauthorized, w)
			return
		}

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Got an error")
			}

			key := settings.GetEnvKey("SECRET_KEY", "MY_SECRET_KEY")
			return []byte(key), nil
		})

		if err != nil {
			helper.Error("Unable to handle request. Token problems. Relogin please and try againg", http.StatusInternalServerError, w)
			return
		}

		if token.Valid {
			next.ServeHTTP(w, r)
		}

	})
}
