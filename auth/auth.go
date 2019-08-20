package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

type tokenRes struct {
	Token string `json:"token"`
}

type tokenRequest struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

// GetTokenHandker get token
var GetTokenHandker = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	if r.Method == http.MethodGet {
		claims["admin"] = true
		claims["sub"] = "123456789"
		claims["name"] = "fuga"
		claims["iat"] = time.Now()
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
		if err != nil {
			w.Write([]byte("err"))
		}
		json.NewEncoder(w).Encode(&tokenRes{Token: tokenString})
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user tokenRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	claims["admin"] = true
	claims["sub"] = user.PassWord
	claims["name"] = user.UserName
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		w.Write([]byte("err"))
	}
	json.NewEncoder(w).Encode(&tokenRes{Token: tokenString})
})

// JwtMiddleware check token
var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})
