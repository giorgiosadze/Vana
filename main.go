package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// User structure to hold user information
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = map[string]string{
	"testuser": "password123",
}

// Function to validate user input
func validateUserInput(user User) error {
	if len(user.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(user.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	return nil
}

func main() {

	db, err := sql.Open("sqlite3", "/data/db/data.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	createQuizzesTable := `
    CREATE TABLE IF NOT EXISTS quizzes (
        quiz_id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT
    );`
	_, err = db.Exec(createQuizzesTable)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("quizzes created successfully")
	}

	createQuestionsTable := `
    CREATE TABLE IF NOT EXISTS questions (
        question_id INTEGER PRIMARY KEY AUTOINCREMENT,
        quiz_id INTEGER,
        question_text TEXT NOT NULL,
        FOREIGN KEY (quiz_id) REFERENCES quizzes(quiz_id)
    );`
	_, err = db.Exec(createQuestionsTable)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("questions created successfully")
	}

	createOptionsTable := `
    CREATE TABLE IF NOT EXISTS options (
        option_id INTEGER PRIMARY KEY AUTOINCREMENT,
        question_id INTEGER,
        option_text TEXT NOT NULL,
        is_correct BOOLEAN,
        FOREIGN KEY (question_id) REFERENCES questions(question_id)
    );`
	_, err = db.Exec(createOptionsTable)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("options created successfully")
	}

	createUsersTable := `
    CREATE TABLE IF NOT EXISTS Users (
        user_id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL,
        password TEXT NOT NULL
    );`
	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("users created successfully")
	}

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/createQ", createQHandler)
	http.HandleFunc("/readQ", readQHandler)
	http.HandleFunc("/updateQ", updateQHandler)
	http.HandleFunc("/deleteQ", deleteQHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))
	log.Println("Listening on port 8000")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if password, ok := users[user.Username]; ok && password == user.Password {
		jsonResponse(w, "Login successful")
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// handle logout
	jsonResponse(w, "Logout endpoint")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = validateUserInput(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create user
	insertUserQuery := `INSERT INTO users (username, password) VALUES (?, ?)`
	_, err = db.Exec(insertUserQuery, user.Username, user.Password)
	if err != nil {
		http.Error(w, "Failed to create user account", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{"message": "User registered successfully"})
}

func createQHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// handle create question
	jsonResponse(w, "CreateQ endpoint")
}

func readQHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// handle read questions
	jsonResponse(w, "ReadQ endpoint")
}

func updateQHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// handle update question
	jsonResponse(w, "UpdateQ endpoint")
}

func deleteQHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// handle delete question
	jsonResponse(w, "DeleteQ endpoint")
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
