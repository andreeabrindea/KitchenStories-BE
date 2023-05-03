package db

type Recipe struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	UserID       string `json:"user_name"`
	Photo        string `json:"photo,omitempty"`
	Ingredients  string `json:"ingredients"`
	Instructions string `json:"instructions"`
}
