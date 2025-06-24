package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrTaskNotFound = errors.New("task not found")
	ErrReferrerNotFound = errors.New("referrer not found")
	ErrTaskAlreadyDone = errors.New("task already done")
	ErrReferrerAlreadyExists = errors.New("referrer already exists")
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}
