package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"task/internal/service"

	"github.com/gorilla/mux"
)

type BodyForCompleteTask struct {
	TaskID string `json:"task_id"`
}

func (h *Handler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	const op = "handler.user.CompleteTask"
	log := h.log.With("op", op)

	idFromToken := r.Header.Get("user_id")
	idFromURL := mux.Vars(r)["id"]
	if idFromToken != idFromURL {
		http.Error(w, "This id is not yours", http.StatusForbidden)
		return
	}
	body, err := io.ReadAll(r.Body)
 	if err != nil {
		http.Error(w, "Something wrong with server", http.StatusInternalServerError)
		log.Error("Failed to read request body", "error", err)
		return
	}
	var taskBody BodyForCompleteTask
	if err := json.Unmarshal(body, &taskBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Error("Failed to unmarshal request body", "error", err)
		return
 	}

	taskID := taskBody.TaskID
	err = h.s.CompleteTask(idFromToken, taskID)
	if err != nil {
		if errors.Is(err, service.ErrUserAndTaskRequired) {
   			http.Error(w, "User ID and Task ID are required", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrInvalidUserID) {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrInvalidTaskID) {
   			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrTaskNotFound) {
   			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrTaskAlreadyDone) {
			http.Error(w ,"task already done", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something wrong with server", http.StatusInternalServerError)
		log.Error("Error occurred with s.CompleteTask ", "error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Info("Task completed successfully", "user_id", idFromToken, "task_id", taskID)
}


type BodyForReferrer struct {
	ReferralID string `json:"referral_id"`
}

func (h *Handler) GetReferrer(w http.ResponseWriter, r *http.Request) {
	const op = "handler.user.GetReferrer"
	log := h.log.With("op", op)

	idFromToken := r.Header.Get("user_id")
	idFromURL := mux.Vars(r)["id"]
	if idFromToken != idFromURL {
		http.Error(w, fmt.Sprintf("This id is not yours. Your ID is %s ",idFromToken), http.StatusForbidden)
		return

	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Something wrong with server", http.StatusInternalServerError)
		log.Error("Failed to read request body", "error", err)
		return
	}
	var referrerBody BodyForReferrer
	if err := json.Unmarshal(body, &referrerBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Error("Failed to unmarshal request body", "error", err)
		return
	}
	referrerID := referrerBody.ReferralID
	if err = h.s.GetReferrer(idFromToken, referrerID); err != nil {
		if errors.Is(err, service.ErrUserAndReferrerRequired) {
			http.Error(w, "User ID and referral ID are required", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrReferrerAlreadyExists) {
			http.Error(w, "Referral already exists", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrReferrerNotFound) {
			http.Error(w, "Referral not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrReferYourself) {
			http.Error(w, "You can't refer yourself", http.StatusBadRequest)
			return
		}
		
		http.Error(w, "Something wrong with server", http.StatusInternalServerError)
		log.Error("Error occurred with s.GetReferrer", "error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Info("Referrer retrieved successfully", "user_id", idFromToken, "referral_id", referrerID)
}