package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"server-kitchen-stories/db"
	"strconv"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	recipes, err := db.GetAllRecipes("postgres://ejyvmpli:6ADd6xq0YUrVCyH0I7s1nfCT1Qv5gMVw@mouse.db.elephantsql.com/ejyvmpli")
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	u, _ := json.MarshalIndent(recipes, "", "  ")
	_, err = w.Write(u)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetRecipesById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := ParseIDFromQueryString(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	recipes, err := db.GetRecipesById("postgres://ejyvmpli:6ADd6xq0YUrVCyH0I7s1nfCT1Qv5gMVw@mouse.db.elephantsql.com/ejyvmpli", id)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	output, _ := json.MarshalIndent(recipes, "", "  ")
	_, err = w.Write(output)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ParseIDFromQueryString(r *http.Request) (int, error) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("id should be an integer")
	}
	return id, nil
}
