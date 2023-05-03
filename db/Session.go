package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

type Session struct {
	UserID       int
	SessionToken string
	Expiry       time.Time
}

func InsertSession(session *Session) error {
	conn, err := pgx.Connect(context.Background(), "postgres://ejyvmpli:6ADd6xq0YUrVCyH0I7s1nfCT1Qv5gMVw@mouse.db.elephantsql.com/ejyvmpli")
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(),
		"INSERT INTO Sessions (user_id, session_token, expiry) VALUES ($1, $2, $3)",
		session.UserID, session.SessionToken, session.Expiry,
	)
	if err != nil {
		return err
	}

	return nil
}
