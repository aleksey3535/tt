package repository

import (
	"database/sql"
	"errors"
)

func (r *Repository) CreateUser(username, password string) (int, error) {
	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id`
	var id int
	err := r.db.QueryRow(query, username, password).Scan(&id)
	if err != nil {
		return 0, ErrUserAlreadyExists
	}
	return id, nil
}

func (r *Repository) GetUser(username string) (int, string, error) {
	query := `SELECT id, password_hash FROM users WHERE login = $1`
	var id int
	var password_hash string
	err := r.db.QueryRow(query, username).Scan(&id, &password_hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", ErrUserNotFound
		}
		return 0, "", err
	}
	return id, password_hash,nil
}