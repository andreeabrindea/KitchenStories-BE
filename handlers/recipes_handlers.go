package handlers

import (
	"encoding/json"
	"net/http"
	"server-kitchen-stories/db"
)

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
	// Parse the JSON request body into a User struct
	var recipe db.Recipe
	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the recipe into the database
	err = db.InsertRecipe(recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated)
}
