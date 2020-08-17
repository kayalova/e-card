package controllers

import (
	"github.com/gorilla/mux"
	"github.com/kayalova/e-card-catalog/middlewares"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	cardsSubrouter := r.PathPrefix("/api/cards").Subrouter()
	authSubrouter := r.PathPrefix("/api/auth").Subrouter()
	booksSubrouter := r.PathPrefix("/api/books").Subrouter()
	schoolsSubrouter := r.PathPrefix("/api/schools").Subrouter()

	// authentication checking
	// http.Handle("/", middleware.IsAuthenticated(cardsSubrouter))

	//auth
	authSubrouter.HandleFunc("/signup", middlewares.SignUp).Methods("POST")
	authSubrouter.HandleFunc("/signin", middlewares.SignIn).Methods("POST")
	//school
	schoolsSubrouter.HandleFunc("/getAll", middlewares.GetAllSchools).Methods("GET")
	//books
	booksSubrouter.HandleFunc("/getAll", middlewares.GetAllBooks).Methods("GET")
	booksSubrouter.HandleFunc("/filter", middlewares.FilterBooks).Methods("GET")
	//cards
	cardsSubrouter.HandleFunc("/create", middlewares.CreateCard).Methods("POST")
	cardsSubrouter.HandleFunc("/filter", middlewares.FilterCards).Methods("GET")
	cardsSubrouter.HandleFunc("/edit/{id}", middlewares.EditCard).Methods("PUT")
	cardsSubrouter.HandleFunc("/getAll", middlewares.GetAllCards).Methods("GET")
	cardsSubrouter.HandleFunc("/getOne/{id}", middlewares.GetOneCard).Methods("GET")
	cardsSubrouter.HandleFunc("/attachBook/{id}", middlewares.AttachToCard).Methods("GET")
	cardsSubrouter.HandleFunc("/detachBook/{id}", middlewares.DetachFromCard).Methods("GET")
	cardsSubrouter.HandleFunc("/delete/{id}", middlewares.DeleteOneCard).Methods("DELETE")

	return r
}
