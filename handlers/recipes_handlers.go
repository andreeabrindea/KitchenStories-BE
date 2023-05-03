package handlers

import (
	"encoding/json"
	"net/http"
	"server-kitchen-stories/db"
)

func GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	enableCors(r, &w)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	recipes, err := db.GetAllRecipes()
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
func AddRecipe(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Set headers for the main request
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Parse the JSON request body into a RecipeAdd struct
	var recipe db.RecipeAdd
	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = db.InsertRecipe(recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated)
}
