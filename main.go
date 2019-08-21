package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/b1018043/jwt_api/auth"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type postTodo struct {
	Todo string `json:"todo"`
}

// Todo is todo
type Todo struct {
	gorm.Model
	Todo    string `json:"todo"`
	Process string `json:"process"`
	User    string `json:"user"`
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
	r.Handle("/todo", todos)
	r.Handle("/private", auth.JwtMiddleware.Handler(usertodos))
	r.Handle("/auth", auth.GetTokenHandker)
	log.Println("port :", *addr)
	if err := http.ListenAndServe(*addr, r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

var todos = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	todo := &Todo{
		Todo:    "ねる",
		Process: "plan",
		User:    "hoge",
	}
	json.NewEncoder(w).Encode(todo)
})

var usertodos = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var todos []Todo
		db.Find(&todos)
		json.NewEncoder(w).Encode(&todos)
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
		db.Create(&Todo{User: "hoge", Todo: posttodo.Todo, Process: "plan"})
	case http.MethodPut:
	case http.MethodDelete:
	}
})
