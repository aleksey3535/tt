package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	const op = "handler.user.GetStatus"
	log := h.log.With("op", op)

	idFromToken := r.Header.Get("user_id")
	idFromURL := mux.Vars(r)["id"]
	if idFromToken != idFromURL {
		http.Error(w, fmt.Sprintf("This id is not yours. Your ID is %s ",idFromToken), http.StatusForbidden)
		return
	}
	userID, err := strconv.Atoi(idFromToken)
	if err != nil {
		http.Error(w, "Something wrong with server", http.StatusInternalServerError)
		log.Error("Failed to convert user_id from string to int", "error", err)
		return
	}

	user, err := h.s.GetUserStatus(userID)
	if err != nil {
		http.Error(w, "Something wrong with server", http.StatusInternalServerError)
		log.Error("Error occurred with s.GetUserStatus", "error", err)
		return
	}
	message, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Something wrong with server", http.StatusInternalServerError)
		log.Error("Failed to marshal user status", "error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(message)
	log.Info("User status retrieved successfully", "user_id", userID)
}

func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	const op = "handler.user.GetLeaderboard"
	log := h.log.With("op", op)

	users, err := h.s.GetLeaderboard()
	if err != nil {
		http.Error(w, "Something wrong with server", http.StatusInternalServerError)
		log.Error("Error occurred with s.GetLeaderboard", "error", err)
		return
	}
	message, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Something wrong with server", http.StatusInternalServerError)
		log.Error("Failed to marshal leaderboard", "error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(message)
	log.Info("Leaderboard retrieved successfully")
}