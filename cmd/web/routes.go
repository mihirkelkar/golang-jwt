package main

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/mihirkelkar/golang-jwt/pkg/models"
)

//Application : Struct to encapsulate important parts of the service.
type Application struct {
	ErrLog      *log.Logger
	InfLog      *log.Logger
	UserService models.UserService
}

func (app *Application) routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/signup", app.SignUp).Methods("POST")
	router.HandleFunc("/login", app.LogIn).Methods("POST")
	router.HandleFunc("/protectedendpoint", TokenAuthenticate(app.ProtectedHandler)).Methods("GET")
	return router
}
