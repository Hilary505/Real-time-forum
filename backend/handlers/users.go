package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
	"real-time-forum/backend/utils"
)

func HandleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	user, err := utils.GetUserFromSession(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}
	// Get all users
	rows, err := database.Db.Query(`
	SELECT id,uuid, gender, age, nickname, email, first_name, last_name, last_message_time
	FROM users
	ORDER BY
	CASE WHEN last_message_time IS NULL THEN 1 ELSE 0 END,
            last_message_time DESC,
	nickname ASC
	`)

	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	// Parse the users
	var users []models.User

	for rows.Next() {
		var u models.User
		var lastMsgAtStr sql.NullString
		err := rows.Scan(&u.ID, &u.UUID, &u.Gender, &u.Age, &u.Nickname, &u.Email, &u.FirstName, &u.LastName, &lastMsgAtStr)
		if err != nil {
			http.Error(w, "Error parsing users", http.StatusInternalServerError)
			return
		}

		if lastMsgAtStr.Valid {
			parsedTime, err := time.Parse(time.RFC3339, lastMsgAtStr.String) // Adjust format if needed
			if err == nil {
				u.LastPrivateMessageAt = parsedTime // Assign to new field in models.User
			} else {
				log.Printf("Error parsing last_private_message_at for user %s: %v", u.Nickname, err)
			}
		}

		if u.ID != user.ID {
			users = append(users, u)
		}
	}

	for i := range users {
		mu.Lock()
		_, ok := Clients[users[i].UUID]
		mu.Unlock()
		if ok {
			users[i].IsOnline = true
		} else {
			users[i].IsOnline = false
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
