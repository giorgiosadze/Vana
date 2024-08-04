package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Username string   `json:"userName"`
	Password string   `json:"passWord"`
	Quizzes  []string `json:"quizzes"`
}

type Response struct {
	UserDetails   User   `json:"userDetails"`
	CurrentQuiz   string `json:"currentQuiz"`
	CurrentAnswer string `json:"currentAnswer"`
}

type Quiz struct {
	Id        int        `json:"Id"`
	Questions []Question `json:"Questions"`
}

type Question struct {
	Id            int      `json:"Id"`
	QuestionText  string   `json:"QuestionText"`
	Answers       []Answer `json:"Answers"`
	CorrectAnswer Answer   `json:"Answer"`
}

type Answer struct {
	Id         int
	AnswerText string
}

func main() {
	var db Data
	db.Connect()
	db.Init()
	http.HandleFunc("/v1/login", loginHandler)
	http.HandleFunc("/v1/logout", logoutHandler)
	http.HandleFunc("/v1/UserRegister", UserRegisterHandler)
	http.HandleFunc("/v1/createQuestion", createQHandler)
	http.HandleFunc("/v1/readQuestion", readQHandler)
	http.HandleFunc("/v1/updateQuestion", updateQHandler)
	http.HandleFunc("/v1/deleteQuestion", deleteQHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
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

	var db Data
	db.Connect()
	u := db.QueryUser(User{user.Username, user.Password, nil})

	if !(u.Username == user.Username && u.Password == user.Password) {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	resp := Response{
		UserDetails: User{
			Username: u.Username,
			Password: u.Password,
			Quizzes:  u.Quizzes,
		},
		CurrentQuiz:   "",
		CurrentAnswer: "",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	//
}

func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	//
}

func createQHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	//
}

func readQHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	//
}

func updateQHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	//
}

func deleteQHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	//
}

type Data struct {
	DB *sql.DB
}

func (d *Data) Connect() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	d.DB = db
}

func (d *Data) Init() {

	createQuizzesTable := `
    CREATE TABLE IF NOT EXISTS quizzes (
        quiz_id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT
    );`
	_, err := d.DB.Exec(createQuizzesTable)
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
	_, err = d.DB.Exec(createQuestionsTable)
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
	_, err = d.DB.Exec(createOptionsTable)
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
	_, err = d.DB.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("users created successfully")
	}
}

// ------- //
// User DB //
// ------- //

func (d *Data) QueryUser(user User) User {

	row := d.DB.QueryRow(`
			SELECT username, password 
			FROM user 
			WHERE username = ? AND password = ?`, user.Username, user.Password)
	var username string
	var password string
	row.Scan(&username, &password)
	return User{username, password, nil}
}
func (d *Data) CreateUser(user User) {
	// TODO
}
func (d *Data) UpdateUser(user User) {
	// TODO
}
func (d *Data) DeleteUser(user User) {
	// TODO
}

// ------- //
// Quiz DB //
// ------- //

func (d *Data) QueryQuiz(quiz Quiz) {
	// TODO
}
func (d *Data) CreateQuiz(quiz Quiz) {
	// TODO
}
func (d *Data) UpdateQuiz(quiz Quiz) {
	// TODO
}
func (d *Data) DeleteQuiz(quiz Quiz) {
	// TODO
}

// --------- //
// Answer DB //
// --------- //

func (d *Data) CreateAnswer(answer Answer) {
	// TODO
}
func (d *Data) UpdateAnswer(answer Answer) {
	// TODO
}
func (d *Data) DeleteAnswer(answer Answer) {
	// TODO
}

// ----------- //
// Question DB //
// ----------- //

func (d *Data) QueryQuestion(answer Answer) {
	// TODO
}
func (d *Data) CreateQuestion(answer Answer) {
	// TODO
}
func (d *Data) UpdateQuestion(answer Answer) {
	// TODO
}
func (d *Data) DeleteQuestion(answer Answer) {
	// TODO
}
