package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"log"
)

func GetAllUsers(connection string) ([]User, error) {
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

	rows, err := conn.Query(context.Background(), "SELECT * FROM Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &user.Country, &user.City)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUsersById(connection string, id int) ([]User, error) {
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

	rows, err := conn.Query(context.Background(), "SELECT * FROM Users WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &user.Country, &user.City)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func InsertUser(user User) error {
	// Open a database connection
	db, err := sql.Open("postgres", "postgres://ejyvmpli:6ADd6xq0YUrVCyH0I7s1nfCT1Qv5gMVw@mouse.db.elephantsql.com/ejyvmpli")
	if err != nil {
		return err
	}
	defer db.Close()

	// Execute the SQL INSERT statement
	stmt, err := db.Prepare("INSERT INTO Users(id, username, email, password, country, city) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.UserName, user.Email, user.Password, user.Country, user.City)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByUsername(username string) (User, error) {
	conn, err := pgx.Connect(context.Background(), "postgres://ejyvmpli:6ADd6xq0YUrVCyH0I7s1nfCT1Qv5gMVw@mouse.db.elephantsql.com/ejyvmpli")
	if err != nil {
		log.Fatal(err)
		return User{}, err
	}
	defer conn.Close(context.Background())

	row, err := conn.Query(context.Background(), "SELECT id, username, email, password, country, city FROM Users WHERE username=$1", username)
	if err != nil {
		return User{}, err
	}
	defer row.Close()

	user := User{}
	if row.Next() {
		err = row.Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &user.Country, &user.City)
		if err != nil {
			return User{}, err
		}
	} else {
		return User{}, fmt.Errorf("user not found")
	}

	return user, nil
}

func CheckPassword(userPassword string, password string) bool {
	return userPassword == password
}
