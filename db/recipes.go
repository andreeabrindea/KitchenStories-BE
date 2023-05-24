package db

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"log"
)

func GetAllRecipes() ([]Recipe, error) {
	conn, err := pgx.Connect(context.Background(), "postgres://ejyvmpli:6ADd6xq0YUrVCyH0I7s1nfCT1Qv5gMVw@mouse.db.elephantsql.com/ejyvmpli")
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

	rows, err := conn.Query(context.Background(), "SELECT r.id, r.name as recipe_name, u.username as user_name, r.photo, r.ingredients, r.instructions FROM Recipe r JOIN Users u ON r.user_id = u.id")
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

	rows, err := conn.Query(context.Background(), "SELECT r.id, r.name as recipe_name, u.username as user_name, r.photo, r.ingredients, r.instructions FROM Recipe r JOIN Users u ON r.user_id = u.id WHERE r.id=$1", id)
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

func InsertRecipe(recipe RecipeAdd) error {
	// Open a database connection
	db, err := sql.Open("postgres", "postgres://ejyvmpli:6ADd6xq0YUrVCyH0I7s1nfCT1Qv5gMVw@mouse.db.elephantsql.com/ejyvmpli")
	if err != nil {
		return err
	}
	defer db.Close()

	// Execute the SQL INSERT statement
	stmt, err := db.Prepare("INSERT INTO Recipe (id, name, user_id, photo, ingredients, instructions) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(recipe.ID, recipe.Name, recipe.UserID, recipe.Photo, recipe.Ingredients, recipe.Instructions)
	if err != nil {
		return err
	}

	return nil
}

func GetRecipesByName(connection string, name string) ([]Recipe, error) {
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

	rows, err := conn.Query(context.Background(), "SELECT r.id, r.name as recipe_name, u.username as user_name, r.photo, r.ingredients, r.instructionsFROM Recipe rJOIN Users u ON r.user_id = u.id WHERE r.name LIKE '%' || $1 || '%'", name)
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
