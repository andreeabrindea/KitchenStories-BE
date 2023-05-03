package db

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Country  string `json:"country"`
	City     string `json:"city"`
}
