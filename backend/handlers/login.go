package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"real-time-forum/backend/database"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Identifier string `json:"identifier"` // nickname or email
	Password   string `json:"password"`
}

type LoginResponse struct {
	ID       string `json:"userid"`
	Nickname string `json:"nickname"`
}

// Handles client login
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var creds LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := database.Db.QueryRow(`SELECT id,uuid, nickname, password FROM users WHERE nickname = ? OR email = ?`, creds.Identifier, creds.Identifier).Scan(&user.ID, &user.UUID, &user.Nickname, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}
	// Create a session
	u, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionID := u.String()

	expiresAt := time.Now().Add(24 * time.Hour) // Sessions expire after 60 sec

	_, err = database.Db.Exec("INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)",
		user.ID, sessionID, expiresAt,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occured while logging you in"))
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Expires:  expiresAt,
		HttpOnly: true,
		Path:     "/",
	})

	response := LoginResponse{
		ID:       user.UUID,
		Nickname: user.Nickname,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
