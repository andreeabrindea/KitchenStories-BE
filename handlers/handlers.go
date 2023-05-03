package handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"server-kitchen-stories/db"
	"strconv"
	"strings"
	"time"
)

func enableCors(r *http.Request, w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")

}

func ParseIDFromPath(r *http.Request) (int, error) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		return 0, errors.New("invalid path")
	}

	idStr := pathParts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("id should be an integer")
	}
	return id, nil
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	enableCors(r, &w)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	users, err := db.GetAllUsers("postgres://ejyvmpli:6ADd6xq0YUrVCyH0I7s1nfCT1Qv5gMVw@mouse.db.elephantsql.com/ejyvmpli")
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	u, _ := json.MarshalIndent(users, "", "  ")
	_, err = w.Write(u)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetUsersById(w http.ResponseWriter, r *http.Request) {
	enableCors(r, &w)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := ParseIDFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	users, err := db.GetUsersByIdFromDB("postgres://ejyvmpli:6ADd6xq0YUrVCyH0I7s1nfCT1Qv5gMVw@mouse.db.elephantsql.com/ejyvmpli", id)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	output, _ := json.MarshalIndent(users, "", "  ")
	_, err = w.Write(output)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}
func CreateUser(w http.ResponseWriter, r *http.Request) {
	enableCors(r, &w)
	// Parse the JSON request body into a User struct
	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the user into the database
	err = db.InsertUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID       int       `json:"user_id"`
	SessionToken string    `json:"session_token"`
	Expiry       time.Time `json:"expiry"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(r, &w)
	// Parse the login request from the request body
	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	username := loginReq.Username
	password := loginReq.Password

	// Validate the credentials against a database
	user, err := db.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if db.CheckPassword(user.Password, password) == false {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create a new session cookie to maintain the authenticated state
	sessionToken, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session := &db.Session{
		UserID:       user.ID,
		SessionToken: sessionToken.String(),
		Expiry:       time.Now().Add(24 * time.Hour),
	}
	err = db.InsertSession(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the session cookie in the response header
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken.String(),
		Expires:  session.Expiry,
		HttpOnly: true,
		Secure:   true,
	})

	// Construct the response struct
	loginResp := LoginResponse{
		UserID:       session.UserID,
		SessionToken: session.SessionToken,
		Expiry:       session.Expiry,
	}

	// Encode the response struct as JSON and write to the response body
	jsonResp, err := json.Marshal(loginResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}
