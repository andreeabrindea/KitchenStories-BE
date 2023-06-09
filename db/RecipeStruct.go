package db

type Recipe struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	UserID       string `json:"user_name"`
	Photo        []byte `json:"photo,omitempty"`
	Ingredients  string `json:"ingredients"`
	Instructions string `json:"instructions"`
}

type RecipeAdd struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	UserID       int    `json:"user_id"`
	Photo        string `json:"photo,omitempty"`
	Ingredients  string `json:"ingredients"`
	Instructions string `json:"instructions"`
}
