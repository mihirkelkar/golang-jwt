package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func TokenAuthenticate(next http.HandlerFunc) http.HandlerFunc {
	fmt.Print("Token Autenticate was called\n")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//This assumes that auth token from generate token is now being passed.
		// in the auth field of the web request.
		var authToken string
		auth := r.Header.Get("Authorization")
		bearerToken := strings.Split(auth, " ")
		if len(bearerToken) == 2 {
			authToken = bearerToken[1]
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Error: Internal Server Error")
				}
				return []byte("secret_key"), nil
			})

			if err != nil {
				fmt.Println(err)
				http.Error(w, "Error: Unauthorized Access from the first error", http.StatusUnauthorized)
				return
			}
			if token.Valid {
				ctx := context.WithValue(r.Context(), "Email", claims.Email)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
		http.Error(w, "Error: Unauthorized Access", http.StatusUnauthorized)
	})
}
