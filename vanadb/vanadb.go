package vanadb

import (
	"database/sql"
	"fmt"
	"log"
	"vana/types"
)

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

func (d *Data) CountUsers() (uint64, error) {
	var count uint64
	query := `SELECT COUNT(*) FROM users`
	err := d.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error querying user count: %w", err)
	}

	return count, nil
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

func (d *Data) QueryUser(user types.User) types.User {

	row := d.DB.QueryRow(`
			SELECT userid, username, password 
			FROM user 
			WHERE username = ? AND password = ?`, user.Username, user.Password)
	var userid uint64
	var username string
	var password string
	row.Scan(&userid, &username, &password)
	return types.User{userid, username, password}
}
