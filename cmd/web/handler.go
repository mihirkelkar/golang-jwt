package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mihirkelkar/golang-jwt/pkg/models"
)

//SignUp : Signs up a person for a service.
func (app *Application) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if user.Email == "" {
		http.Error(w, "Error: Email is Missing", 500)
		return
	}

	if user.Password == "" {
		http.Error(w, "Error: Password is Missing", 500)
		return
	}

	err = app.UserService.Insert(&user)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	return
}

//LogIn : Logs the function in for a service.
func (app *Application) LogIn(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var jwt models.JWT
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error: Invalid Request", http.StatusBadRequest)
	}
	if user.Email == "" {
		http.Error(w, "Error: Email missing", http.StatusBadRequest)
		return
	}

	if user.Password == "" {
		http.Error(w, "Error: Password Missing", http.StatusBadRequest)
		return
	}

	err = app.UserService.Authenticate(&user)
	if err != nil {
		http.Error(w, "Error: Error Autenticating User", http.StatusInternalServerError)
		return
	}

	token, err := app.GenerateToken(&user, "secret_key")
	if err != nil {
		http.Error(w, "Error: Error Authenticating User", http.StatusInternalServerError)
		return
	}
	jwt.Token = token
	w.Write([]byte(jwt.Token))
	return
}

//protectedHandler : End point that is protected by middleware.
func (app *Application) ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Protected End Point Invoked\n")
	email := r.Context().Value("Email")
	foundUser, err := app.UserService.ByEmail(email.(string))
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Print(foundUser)
	return
}
