package handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"server-kitchen-stories/db"
	"strconv"
	"time"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func ParseIDFromQueryString(r *http.Request) (int, error) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("id should be an integer")
	}
	return id, nil
}
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := ParseIDFromQueryString(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	users, err := db.GetUsersById("postgres://ejyvmpli:6ADd6xq0YUrVCyH0I7s1nfCT1Qv5gMVw@mouse.db.elephantsql.com/ejyvmpli", id)
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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

	loginResp := db.Session{
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
	_, err = w.Write(jsonResp)
	if err != nil {
		return
	}

}
