package repository

import (
	"database/sql"
	"errors"
	"task/internal/models"
)

func (r *Repository) GetUserStatus(id int) (models.UserForStatus, error) {
	query := `SELECT id, login, quantity, points, referrer_id FROM users WHERE id = $1`
	var user models.UserForStatus
	err := r.db.QueryRow(query, id).Scan(&user.Id, &user.Login, &user.Quantity, &user.Points, &user.ReferrerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserForStatus{}, ErrUserNotFound
		}
		return models.UserForStatus{}, errors.New("with db.QueryRow" + err.Error())
	}
	query = `SELECT task_id FROM completed_tasks WHERE user_id = $1`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return models.UserForStatus{}, errors.New("with db.Query " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var taskID int
		if err = rows.Scan(&taskID); err != nil {
			return models.UserForStatus{}, errors.New("with rows.Scan " + err.Error())
		}
		user.CompletedTasks = append(user.CompletedTasks, taskID)
	}
	return user, nil
}

func (r *Repository) GetLeaderboard() ([]models.UserForLeaderBoard, error) {
	query := `SELECT login, quantity, points FROM users ORDER BY points DESC LIMIT 3`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, errors.New("with db.Query" + err.Error())
	}
	defer rows.Close()

	var users []models.UserForLeaderBoard
	for rows.Next() {
		var user models.UserForLeaderBoard
		if err := rows.Scan(&user.Login, &user.Quantity, &user.Points); err != nil {
			return nil, errors.New("with rows.Scan" + err.Error())
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *Repository) CompleteTask(userID, taskID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		return errors.New("with db.Begin" + err.Error())
	}
	defer tx.Commit()

	query := `SELECT given_points FROM tasks WHERE id = $1` // получаем данные о количестве очков за выполнение задания
	var givenPoints int
	if err = tx.QueryRow(query, taskID).Scan(&givenPoints); err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTaskNotFound
		}
		return errors.New("with tx.QueryRow" + err.Error())
	}
	query = `SELECT task_id FROM completed_tasks WHERE user_id = $1` // проверяем была ли задача выполнена ранее
	rows ,err := tx.Query(query, userID)
	if err != nil {
		tx.Rollback()
		return errors.New("with rows.Query " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var expectedTaskID int
		if err := rows.Scan(&expectedTaskID); err != nil {
			tx.Rollback()
			return errors.New("with rows.Scan " + err.Error())
		}
		if expectedTaskID == taskID {
			tx.Rollback()
			return ErrTaskAlreadyDone
		}
	}
	query = `UPDATE users SET quantity = quantity + 1, points = points + $1 WHERE id = $2` // добавляем пользователю баллы за задание
	_, err = tx.Exec(query, givenPoints, userID)
	if err != nil {
		tx.Rollback()
		return errors.New("with db.Exec update users" + err.Error())
	}
	query = `INSERT INTO completed_tasks (user_id, task_id) VALUES ($1, $2)` // ведем список завершенных заданий для пользователя
	if _, err = tx.Exec(query, userID, taskID); err != nil {
		tx.Rollback()
		return errors.New("with db.Exec completed tasks " + err.Error())
	}
	return nil
}

func (r *Repository) GetReferrer(userID, referrerID int) error {
	query := `SELECT id FROM users WHERE id = $1` // проверка сущестовавания реферала
	var id int
	tx, err :=  r.db.Begin()
 	if err != nil {
		tx.Rollback()
		return errors.New("with db.Begin" + err.Error())
 	}
	defer tx.Commit()
	if err = tx.QueryRow(query, referrerID).Scan(&id); err != nil {	
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			   return ErrReferrerNotFound
		}
  	  	return errors.New("with tx.QueryRow" + err.Error())
	}
	query = `SELECT referrer_id FROM users WHERE id = $1` // проверка наличия реферала у пользователя
	var expectedReferrerID int
	if err = tx.QueryRow(query, userID).Scan(&expectedReferrerID); !errors.Is(err, sql.ErrNoRows) && expectedReferrerID != 0 {
		tx.Rollback()
		if err == nil {
			return ErrReferrerAlreadyExists
		}
		return errors.New("with tx.QueryRow select " + err.Error())
	}
	query = `UPDATE users SET referrer_id = $1 WHERE id = $2` // установка реферала пользователю
	if _,err = tx.Exec(query, referrerID, userID); err != nil {
		return errors.New("with tx.Exec " + err.Error())
	}
	return nil
}
