package controller

import (
	"github.com/gorilla/mux"
	"github.com/kayalova/e-card-catalog/handler"
	"github.com/kayalova/e-card-catalog/middleware"
)

// Router ...
func Router() *mux.Router {
	r := mux.NewRouter()

	cardsSubrouter := r.PathPrefix("/api/cards").Subrouter()
	authSubrouter := r.PathPrefix("/api/auth").Subrouter()
	booksSubrouter := r.PathPrefix("/api/books").Subrouter()
	schoolsSubrouter := r.PathPrefix("/api/schools").Subrouter()

	// authentication check
	cardsSubrouter.Use(middleware.IsAuthorized)
	booksSubrouter.Use(middleware.IsAuthorized)
	schoolsSubrouter.Use(middleware.IsAuthorized)

	//auth
	authSubrouter.HandleFunc("/signup", handler.SignUp).Methods("POST")
	authSubrouter.HandleFunc("/signin", handler.SignIn).Methods("POST")
	//school
	schoolsSubrouter.HandleFunc("/getAll", handler.GetAllSchools).Methods("GET")
	//books
	booksSubrouter.HandleFunc("/getAll", handler.GetAllBooks).Methods("GET")
	booksSubrouter.HandleFunc("/filter", handler.FilterBooks).Methods("GET")
	//cards
	cardsSubrouter.HandleFunc("/create", handler.CreateCard).Methods("POST")
	cardsSubrouter.HandleFunc("/filter", handler.FilterCards).Methods("GET")
	cardsSubrouter.HandleFunc("/edit/{id}", handler.EditCard).Methods("PUT")
	cardsSubrouter.HandleFunc("/getAll", handler.GetAllCards).Methods("GET")
	cardsSubrouter.HandleFunc("/getOne/{id}", handler.GetOneCard).Methods("GET")
	cardsSubrouter.HandleFunc("/attachBook/{id}", handler.AttachToCard).Methods("GET")
	cardsSubrouter.HandleFunc("/detachBook/{id}", handler.DetachFromCard).Methods("GET")
	cardsSubrouter.HandleFunc("/delete/{id}", handler.DeleteOneCard).Methods("DELETE")

	return r
}
