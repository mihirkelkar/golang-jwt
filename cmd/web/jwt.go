package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mihirkelkar/golang-jwt/pkg/models"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

//Generates the JWT Token to be used here.
func (app *Application) GenerateToken(user *models.User, secretkey string) (string, error) {
	expirationTime := time.Now().Add(25 * time.Minute)
	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	stringtoken, err := token.SignedString([]byte(secretkey))
	if err != nil {
		app.ErrLog.Fatal(err)
	}
	return stringtoken, nil
}
