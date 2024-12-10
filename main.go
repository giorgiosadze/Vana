package main

import (
	_ "database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"vana/types"
	"vana/vanadb"
	_ "vana/vanadb"
)

func main() {
	var db vanadb.Data
	db.Connect()
	db.Init()
	http.HandleFunc("POST /v1/login", LoginHandler)
	http.HandleFunc("GET /v1/UserCount", UserCount)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func UserCount(w http.ResponseWriter, r *http.Request) {

	var db vanadb.Data
	db.Connect()

	c, err := db.CountUsers()
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user types.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var db vanadb.Data
	db.Connect()
	u := db.QueryUser(types.User{
		Username: user.Username,
		Password: user.Password})

	if !(u.Username == user.Username && u.Password == user.Password) {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Test")

}
