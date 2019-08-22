package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"github.com/b1018043/jwt_api/database"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

const during = 24

type tokenRes struct {
	Token string `json:"token"`
}

type tokenRequest struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

type loginRequest struct {
	Email string `json:"email"`
	Pass  string `json:"password"`
}

type signUpRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	PassWord string `json:"password"`
}

// DispatchToken return tokenstring
func DispatchToken(sub, name, secret string, t time.Time) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["sub"] = sub
	claims["name"] = name
	claims["iat"] = t
	claims["exp"] = t.Add(time.Hour * during).Unix()
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

// JwtMiddleware check token
var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

// LoginHandler is http.HandlerFunc
var LoginHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var tmp loginRequest
	if err := json.NewDecoder(r.Body).Decode(&tmp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var user database.User
	if err := database.GetDB().Where("email=? AND password=?", tmp.Email, tmp.Pass).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	tokenString, err := DispatchToken(user.UserID, user.UserName, os.Getenv("SECRET_KEY"), time.Now())
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	json.NewEncoder(w).Encode(&tokenRes{Token: tokenString})
})

// SignUpHandler use when signup
var SignUpHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var tmp signUpRequest
	if err := json.NewDecoder(r.Body).Decode(&tmp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user database.User
	if err := database.GetDB().Where("email=?", tmp.Email).First(&user).Error; !gorm.IsRecordNotFoundError(err) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := uuid.NewRandom()
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	var setuser = database.User{UserID: u.String(), UserName: tmp.UserName, Password: tmp.PassWord, Email: tmp.Email}
	if err := database.GetDB().Create(&setuser).Error; err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	tokenString, err := DispatchToken(setuser.UserID, setuser.UserName, os.Getenv("SECRET_KEY"), time.Now())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(&tokenRes{Token: tokenString})
})
