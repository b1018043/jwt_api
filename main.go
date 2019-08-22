package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/b1018043/jwt_api/auth"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type postTodo struct {
	Todo string `json:"todo"`
}

// ResponseJSON is struct
type ResponseJSON struct {
	Todos  []Todo `json:"todos"`
	Length int    `json:"length"`
}

// Todo is struct about todo information
type Todo struct {
	gorm.Model
	Todo    string `json:"todo"`
	Process string `json:"process"`
	UserID  string `json:"userid"`
	TodoID  string `json:"todoid"`
}

// User is user information
type User struct {
	gorm.Model
	UserName string `json:"username"`
	UserID   string `json:"userid"`
	Password string `json:"password"`
	Email    string `json:"Email"`
}

func envLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var db *gorm.DB
var er error

func init() {
	db, er = gorm.Open("sqlite3", "./data.db")
	if er != nil {
		return
	}
	db.AutoMigrate(&Todo{})
}

func main() {
	envLoad()
	defer db.Close()
	var addr = flag.String("addr", ":8080", "address")
	flag.Parse()
	r := mux.NewRouter()
	r.Handle("/private", auth.JwtMiddleware.Handler(usertodos))
	r.Handle("/auth", auth.GetTokenHandker)
	log.Println("port :", *addr)
	if err := http.ListenAndServe(*addr, r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

var usertodos = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var todos []Todo
		db.Find(&todos)
		json.NewEncoder(w).Encode(&ResponseJSON{Todos: todos, Length: len(todos)})
	case http.MethodPost:
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var posttodo postTodo
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&posttodo); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		claims := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
		userid, ok := claims["sub"].(string)
		if !ok {
			w.WriteHeader(http.StatusNonAuthoritativeInfo)
			return
		}
		u, err := uuid.NewRandom()
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			return
		}
		db.Create(&Todo{UserID: userid, Todo: posttodo.Todo, Process: "plan", TodoID: u.String()})
	case http.MethodPut:
	case http.MethodDelete:
	}
})
