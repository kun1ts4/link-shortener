package handlers

import (
	"encoding/json"
	"link-shortener/internal/config"
	"link-shortener/internal/domain"
	"log/slog"
	"net/http"
)

type Handler struct {
	uc  domain.UseCase
	log *slog.Logger
	cfg *config.RateLimitConfig
}

func NewHandler(uc domain.UseCase, log *slog.Logger, cfg *config.RateLimitConfig) *Handler {
	return &Handler{
		uc:  uc,
		log: log,
		cfg: cfg,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	b, _ := json.Marshal(errorResponse{Error: msg})
	_, _ = w.Write(b)
}
