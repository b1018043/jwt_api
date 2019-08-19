package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type todo struct {
	Todo    string `json:"todo"`
	Process string `json:"process"`
	User    string `json:"user"`
}

func main() {
	r := mux.NewRouter()
	r.Handle("/todo", todos)
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
