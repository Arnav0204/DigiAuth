package auth

import (
	"context"
	"digiauth/database"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
}

// GenerateJWT generates a JWT token for a given user
func GenerateJWT(email string, role string, id uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 2 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(os.Getenv("SECRET_KEY"))
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	email := req.Email
	password := req.Password
	username := req.Username
	role := req.Role

	if email == "" || password == "" || username == "" || role == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hash the password to store in database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Insert the user into the database
	query := `INSERT INTO users (email, username, password, role) VALUES ($1, $2, $3, $4)`
	_, err = database.DB.Exec(context.Background(), query, req.Email, req.Username, hashedPassword, req.Role)
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered successfully"))
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	email := req.Email
	password := req.Password

	if email == "" || password == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Query to retrieve the complete row
	var user User
	query := "SELECT id, email, username, password, role FROM users WHERE email=$1"
	err := database.DB.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database query error", http.StatusInternalServerError)
		}
		return
	}

	// Compare the provided password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := GenerateJWT(user.Email, user.Role, user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Return the token in the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}
