package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kayalova/e-card/middleware"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	cardsSubrouter := r.PathPrefix("/api/cards").Subrouter()
	authSubrouter := r.PathPrefix("/api/auth").Subrouter()

	// authentication checking
	http.Handle("/", middleware.IsAuthenticated(cardsSubrouter))

	//auth
	authSubrouter.HandleFunc("/signup", middleware.SignUp).Methods("POST")
	authSubrouter.HandleFunc("/signin", middleware.SignIn).Methods("POST")
	//cards
	cardsSubrouter.HandleFunc("/filter", middleware.FilterCards).Methods("GET")
	cardsSubrouter.HandleFunc("/edit/{id}", middleware.EditCard).Methods("PUT")
	cardsSubrouter.HandleFunc("/getAll", middleware.getAllCards).Methods("GET")
	cardsSubrouter.HandleFunc("/getOne/{id}", middleware.getOneCard).Methods("GET")
	cardsSubrouter.HandleFunc("/delete/{id}", middleware.deleteOne).Methods("DELETE")

}

/*
/api/cards/edit/12312

*/
