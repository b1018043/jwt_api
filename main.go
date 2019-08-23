package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/b1018043/jwt_api/database"

	"github.com/b1018043/jwt_api/auth"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type postTodo struct {
	Todo string `json:"todo"`
}

// ResponseJSON is struct
type ResponseJSON struct {
	Todos  []database.Todo `json:"todos"`
	Length int             `json:"length"`
}

type updateTodo struct {
	TodoID  string `json:"todoid"`
	Process string `json:"process"`
}

type deleteTodo struct {
	TodoID string `json:"todoid"`
}

func envLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	envLoad()
	var addr = flag.String("addr", ":8080", "address")
	flag.Parse()
	r := mux.NewRouter()
	r.Handle("/private", auth.JwtMiddleware.Handler(usertodos))
	r.Handle("/login", auth.LoginHandler)
	r.Handle("/signup", auth.SignUpHandler)
	log.Println("port ", *addr)
	if err := http.ListenAndServe(*addr, r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userid, ok := claims["sub"].(string)
	if !ok {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	var todos []database.Todo
	database.GetDB().Where("user_id=?", userid).Find(&todos)
	json.NewEncoder(w).Encode(&ResponseJSON{Todos: todos, Length: len(todos)})
}

var usertodos = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTodo(w, r)
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
		database.GetDB().Create(&database.Todo{UserID: userid, Todo: posttodo.Todo, Process: "plan", TodoID: u.String()})
		getTodo(w, r)
	case http.MethodPatch:
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		claims := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
		userid, ok := claims["sub"].(string)
		if !ok {
			w.WriteHeader(http.StatusExpectationFailed)
			return
		}
		var info updateTodo
		if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		var proc string
		switch info.Process {
		case "plan":
			proc = "doing"
		case "doing":
			proc = "done"
		case "done":
			database.GetDB().Delete(&database.Todo{}, "user_id=? AND todo_id=?", userid, info.TodoID)
			return
		default:
			proc = "plan"
		}
		if err := database.GetDB().Model(&database.Todo{}).Where("user_id=? AND todo_id=?", userid, info.TodoID).Update("process", proc).Error; err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			return
		}
		getTodo(w, r)
	case http.MethodDelete:
		claims := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
		userid, ok := claims["sub"].(string)
		if !ok {
			w.WriteHeader(http.StatusExpectationFailed)
			return
		}
		var info deleteTodo
		if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := database.GetDB().Delete(&database.Todo{}, "user_id=? AND todo_id=?", userid, info.TodoID).Error; err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		getTodo(w, r)
	}
})
