package handler

import (
	"log/slog"
	"net/http"
	"task/internal/config"
	"task/internal/middleware"
	"task/internal/models"

	"github.com/gorilla/mux"
)

type ServiceI interface {
	Register(username, password string) (int, error)
	Login(username, password string) (string, error)
	GetUserStatus(id int) (models.UserForStatus, error)
	GetLeaderboard() ([]models.UserForLeaderBoard, error)
	CompleteTask(userID, taskID string) error
	GetReferrer(userID, referrerID string) error
}

type Handler struct {
	cfg *config.Config
	log *slog.Logger
	s   ServiceI
	mw *middleware.MiddleWare
}

func New(log *slog.Logger, s ServiceI, cfg *config.Config, mw *middleware.MiddleWare) *Handler {
	return &Handler{
		log: log,
		s:   s,
		cfg: cfg,
		mw:  mw,
	}
}

func (h *Handler) InitRoutes() *mux.Router{
	mux := mux.NewRouter()
	mux.Use(h.mw.UseHeader)
	mux.HandleFunc("/users/register", h.Register).Methods(http.MethodPost)
	mux.HandleFunc("/users/login", h.Login).Methods(http.MethodPost)
	mux.HandleFunc("/users/{id}/status", h.mw.CheckAuth(h.GetStatus)).Methods(http.MethodGet)
	mux.HandleFunc("/users/leaderboard", h.mw.CheckAuth(h.GetLeaderboard)).Methods(http.MethodGet)
	mux.HandleFunc("/users/{id}/task/complete", h.mw.CheckAuth(h.CompleteTask)).Methods(http.MethodPost)
	mux.HandleFunc("/users/{id}/referrer", h.mw.CheckAuth(h.GetReferrer)).Methods(http.MethodPost)
	return mux
}
