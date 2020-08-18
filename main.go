package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kayalova/e-card-catalog/controllers"
	_ "github.com/lib/pq"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	r := controllers.Router()
	port := os.Getenv("SERVER_PORT")
	http.ListenAndServe(":"+port, r)

}
