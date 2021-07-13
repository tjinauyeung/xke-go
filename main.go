package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db := initDB()

	r := mux.NewRouter()
	h := &Handler{db: db}

	r.HandleFunc("/users", h.getUsers).Methods("GET")
	r.HandleFunc("/users", h.createUser).Methods("POST")
	r.HandleFunc("/users/{id}", h.getUser).Methods("GET")

	log.Println("Server started: running on port 3000")
	log.Fatal(http.ListenAndServe("localhost:3000", r))
}

type Handler struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

// fill in your db config here
const (
	host     = "localhost"
	user     = "postgres"
	password = "password"
	dbname   = "golangdb"
	port     = 5432
)

func initDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{})

	return db
}

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	var uu []User
	tx := h.db.Find(&uu)
	if tx.Error != nil {
		respond(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, uu, http.StatusOK)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	json.NewDecoder(r.Body).Decode(&u)

	v := validator.New()
	err := v.Struct(u)
	if err != nil {
		respond(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := h.db.Create(&u)
	if tx.Error != nil {
		respond(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, u, http.StatusCreated)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	uid, err := strconv.Atoi(id)
	if err != nil {
		respond(w, err.Error(), http.StatusBadRequest)
		return
	}

	var u User
	tx := h.db.First(&u, uid)
	if tx.Error != nil {
		respond(w, tx.Error.Error(), http.StatusNotFound)
		return
	}

	respond(w, u, http.StatusOK)
}

func respond(w http.ResponseWriter, body interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}
