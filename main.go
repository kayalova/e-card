package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/kayalova/e-card-catalog/controller"
	"github.com/kayalova/e-card-catalog/settings"
	_ "github.com/lib/pq"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	r := controller.Router()
	port := settings.GetEnvKey("SERVER_PORT", "8080")
	http.ListenAndServe(":"+port, r)
}
