package service

import (
	"errors"
	"task/internal/config"
	"task/internal/models"
)


var (
	ErrWrongPassword = errors.New("wrong password")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrBadCredentials    = errors.New("bad credentials")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidUserID	= errors.New("invalid user ID")
	ErrInvalidTaskID	= errors.New("invalid task ID")
	ErrTaskNotFound	= errors.New("task not found")
	ErrUserAndTaskRequired = errors.New("user ID and task ID cannot be empty")
	ErrUserAndReferrerRequired = errors.New("user ID and referrer ID cannot be empty")
	ErrInvalidReferrerID = errors.New("invalid referrer ID")
	ErrTaskAlreadyDone = errors.New("task already done")
	ErrReferrerNotFound = errors.New("referrer not found")
	ErrReferrerAlreadyExists = errors.New("referrer already exists")
	ErrReferYourself = errors.New("you cant refer yourself")
)

type RepositoryI interface {
	CreateUser(username, password string) (int, error)
	GetUser(username string) (int, string, error)
	GetUserStatus(id int) (models.UserForStatus, error)
	GetLeaderboard() ([]models.UserForLeaderBoard, error)
	CompleteTask(userID, taskID int) error
	GetReferrer(userID, referrerID int) error
}

type Service struct {
	repo RepositoryI
	cfg  *config.Config
}

func New(repo RepositoryI, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}
