package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"task/internal/service"
)

type InputBody struct {
	Login    string `json:"login"`
 	Password string `json:"password"`
}

type OutputBody struct {
	Id    int    `json:"id"`
	
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "handler.auth.Register"
	log := h.log.With("op", op)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Something wrong with server", http.StatusInternalServerError)
		log.Error("Failed to read request body", "error", err)
		return
	}
	var input InputBody
	 if err := json.Unmarshal(body, &input); err != nil {
	    http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	id, err := h.s.Register(input.Login, input.Password)
	if err != nil {
		if  errors.Is(err, service.ErrUserAlreadyExists) {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		if errors.Is(err, service.ErrBadCredentials) {
  
			http.Error(w, "Bad credentials", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something wrong with server", http.StatusInternalServerError)
		log.Error("Failed to register user", "error", err)
		return
	}

	output := OutputBody{ Id:  id }
	message, err  := json.Marshal(output)
	if err != nil {
		http.Error(w, "Something wrong with server", http.StatusInternalServerError)
		log.Error("Failed to marshal output body " + err.Error())
		return
	}
	log.Info("User registered successfully", "user_id", id)
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

type OutputForLogin struct {
 Token string `json:"token"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handler.auth.Login"
	log := h.log.With("op", op)
	body, err := io.ReadAll(r.Body)
	 if err != nil { 
   		http.Error(w, "Something wrong with server", http.StatusInternalServerError)
   		log.Error("Failed to read request body", "error", err)
		return
	}
	var input InputBody
	if err := json.Unmarshal(body, &input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
  		return
	}
	token, err := h.s.Login(input.Login, input.Password)
	 if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
  		}
		if errors.Is(err, service.ErrBadCredentials) {
  
			http.Error(w, "Bad credentials", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrWrongPassword) {
			http.Error(w, "Wrong password", http.StatusBadRequest)
			return
		}
		http.Error(w, "Something wrong with server", http.StatusInternalServerError)
		log.Error("Failed to login user", "error", err)
		return
	}
	output := OutputForLogin{Token: token}
	message, err := json.Marshal(output)
	if err != nil {
		http.Error(w, "Something wrong with server", http.StatusInternalServerError)
		log.Error("Failed to marshal output body", "error", err)
		return
	}
	log.Info("User logged in successfully", "token", token)
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}


