package middleware

import (
	"log/slog"
	"net/http"
	"strconv"
	"task/internal/service"
)

type MiddleWare struct {
	log *slog.Logger
	
	s   *service.Service
}

func New(log *slog.Logger, s *service.Service) *MiddleWare {
	return &MiddleWare{
		log: log,
		s:   s,
	}
}

func (m *MiddleWare) UseHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") 
		m.log.Info("Request received", "method", r.Method, "url", r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func (m *MiddleWare) CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Token")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		id, err := m.s.ParseToken(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		r.Header.Set("user_id", strconv.Itoa(id))
		next.ServeHTTP(w, r)
	})	
}