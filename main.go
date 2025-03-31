package main

import (
	_ "database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"vana/types"
	"vana/vanadb"
	_ "vana/vanadb"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
)

// Global database connection
var database vanadb.Data

// Response structures for consistent API responses
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type CountResponse struct {
	Count uint64 `json:"count"`
}

type LoginResponse struct {
	User  types.User `json:"user"`
	Token string     `json:"token"`
}

func init() {
	// Initialize database connection once at startup
	database.Connect()
	database.Init()
}

func main() {
	// Routes
	http.HandleFunc("POST /v1/login", LoginHandler)
	http.HandleFunc("GET /v1/UserCount", UserCount)
	http.HandleFunc("GET /v1/quizzes", AuthMiddleware(GetQuizzes))

	// Serve static files for the frontend
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println("Server starting on :8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// UserCount returns the total number of users
func UserCount(w http.ResponseWriter, r *http.Request) {
	c, err := database.CountUsers()
	if err != nil {
		// Return error as JSON response instead of crashing
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ApiResponse{
			Success: false,
			Message: "Failed to retrieve user count",
		})
		log.Printf("Error counting users: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Success: true,
		Data:    CountResponse{Count: c},
	})
}

// LoginHandler authenticates a user and returns a JWT token
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user types.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Validate input
	if user.Username == "" || user.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiResponse{
			Success: false,
			Message: "Username and password are required",
		})
		return
	}

	// Query for user
	u := database.QueryUser(types.User{
		Username: user.Username,
		Password: user.Password,
	})

	// Check credentials
	if u.UserID == 0 || u.Username != user.Username {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ApiResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  u.UserID,
		"username": u.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your_secret_key")) // Use env var in production!
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ApiResponse{
			Success: false,
			Message: "Failed to generate token",
		})
		return
	}

	// Sanitize user response (don't return password)
	u.Password = ""

	// Return success response with token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Success: true,
		Message: "Login successful",
		Data: LoginResponse{
			User:  u,
			Token: tokenString,
		},
	})
}

// AuthMiddleware verifies the JWT token
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ApiResponse{
				Success: false,
				Message: "Authorization header required",
			})
			return
		}

		// Extract and validate token
		tokenString := authHeader[7:] // Remove "Bearer " prefix
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_secret_key"), nil // Use env var in production!
		})

		if err != nil || !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ApiResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
			return
		}

		// Call the next handler
		next(w, r)
	}
}

// GetQuizzes returns all available quizzes
func GetQuizzes(w http.ResponseWriter, r *http.Request) {
	// This would be implemented to fetch quizzes from the database
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Success: true,
		Data:    []types.Quiz{},
	})
}
