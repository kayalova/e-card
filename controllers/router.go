package controllers

import (
	"github.com/gorilla/mux"
	"github.com/kayalova/e-card-catalog/middleware"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	cardsSubrouter := r.PathPrefix("/api/cards").Subrouter()
	authSubrouter := r.PathPrefix("/api/auth").Subrouter()

	// authentication checking
	// http.Handle("/", middleware.IsAuthenticated(cardsSubrouter))

	//auth
	authSubrouter.HandleFunc("/signup", middleware.SignUp).Methods("POST")
	authSubrouter.HandleFunc("/signin", middleware.SignIn).Methods("POST")
	//cards
	cardsSubrouter.HandleFunc("/create", middleware.Create).Methods("POST")
	cardsSubrouter.HandleFunc("/filter", middleware.Filter).Methods("GET")
	cardsSubrouter.HandleFunc("/edit/{id}", middleware.Edit).Methods("PUT")
	cardsSubrouter.HandleFunc("/getAll", middleware.GetAll).Methods("GET")
	cardsSubrouter.HandleFunc("/getOne/{id}", middleware.GetOne).Methods("GET")
	cardsSubrouter.HandleFunc("/delete/{id}", middleware.DeleteOne).Methods("DELETE")

	return r
}
