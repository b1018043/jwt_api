package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/b1018043/jwt_api/auth"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type todo struct {
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

func main() {
	envLoad()
	r := mux.NewRouter()
	r.Handle("/todo", todos)
	r.Handle("/private", auth.JwtMiddleware.Handler(privateTodo))
	r.Handle("/auth", auth.GetTokenHandker)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

var todos = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	todo := &todo{
		Todo:    "ねる",
		Process: "plan",
		User:    "hoge",
	}
	json.NewEncoder(w).Encode(todo)
})

var privateTodo = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	todo := &todo{
		Todo:    "aaa",
		Process: "bbb",
		User:    "ccc",
	}
	json.NewEncoder(w).Encode(todo)
})
