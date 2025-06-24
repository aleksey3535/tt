package service

import (
	"errors"
	"strconv"
	"task/internal/models"
	"task/internal/repository"
)

func (s *Service) GetUserStatus(id int) (models.UserForStatus, error) {
	user, err := s.repo.GetUserStatus(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return models.UserForStatus{}, ErrUserNotFound
		}
		return models.UserForStatus{}, err
	}
	
	return user, nil
}

func (s *Service) GetLeaderboard() ([]models.UserForLeaderBoard, error) {
 	return s.repo.GetLeaderboard()
}

func (s *Service) CompleteTask(userID, taskID string) error {
	if userID == "" || taskID == "" {
		return ErrUserAndTaskRequired
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return ErrInvalidUserID
	}
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		return ErrInvalidTaskID
	}
	err = s.repo.CompleteTask(userIDInt, taskIDInt)
	if err != nil {
		if  errors.Is(err, repository.ErrTaskNotFound) {
			return ErrTaskNotFound
		}
		if errors.Is(err, repository.ErrTaskAlreadyDone) {
			return ErrTaskAlreadyDone
		}
		return err
	}
	return nil
}

func (s *Service) GetReferrer(userID, referrerID string)  error {
	if userID == "" || referrerID == "" {
		return ErrUserAndReferrerRequired
  	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
   		return ErrInvalidUserID
  	}
	referrerIDInt, err := strconv.Atoi(referrerID)
	if err != nil {
	 	return ErrInvalidReferrerID
   	}
	if userIDInt == referrerIDInt {
		return ErrReferYourself
	}
	if err = s.repo.GetReferrer(userIDInt, referrerIDInt); err != nil {
		if errors.Is(err, repository.ErrReferrerNotFound) {
			return ErrReferrerNotFound
		}
		if errors.Is(err, repository.ErrReferrerAlreadyExists) {
			return ErrReferrerAlreadyExists
		}
		return err
	}
	return nil
}