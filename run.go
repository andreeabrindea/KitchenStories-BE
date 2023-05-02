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
		"/api/getRecipes/{id}",
		handlers.GetAllRecipes,
	)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
