package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

func GetAllRecipes(connection string) ([]Recipe, error) {
	conn, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(conn, context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM Recipe")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		recipe := Recipe{}
		err = rows.Scan(&recipe.ID, &recipe.Name, &recipe.UserID, &recipe.Photo, &recipe.Ingredients, &recipe.Instructions)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

func GetRecipesById(connection string, id int) ([]Recipe, error) {
	conn, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(conn, context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM Recipe WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		recipe := Recipe{}
		err = rows.Scan(&recipe.ID, &recipe.Name, &recipe.UserID, &recipe.Photo, &recipe.Ingredients, &recipe.Instructions)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return recipes, nil
}
