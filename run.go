package main

import (
	"fmt"
	"net/http"
	"server-kitchen-stories/handlers"
)

func main() {

	http.HandleFunc(
		"/api/getRecipes",
		handlers.GetAllRecipes,
	)
	http.HandleFunc(
		"/api/getRecipes/",
		handlers.GetRecipesById,
	)
	http.HandleFunc(
		"/api/getUsers",
		handlers.GetAllUsers,
	)
	http.HandleFunc(
		"/api/getUsers/",
		handlers.GetUsersById,
	)
	http.HandleFunc(
		"/api/users",
		handlers.CreateUser,
	)

	http.HandleFunc(
		"/api/recipes",
		handlers.AddRecipe,
	)
	http.HandleFunc(
		"/api/recipes/",
		handlers.GetRecipesByName)

	http.HandleFunc(
		"/login",
		handlers.LoginHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
