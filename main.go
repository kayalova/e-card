package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	//	_ "github.com/lib/pq" 
	//	"database/sql"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	port := os.Getenv("SERVER_PORT")
	http.ListenAndServe(":"+port, nil)
}
