// Package handlers implements HTTP request handlers for the link-shortener API.
package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/kun1ts4/link-shortener/internal/config"
	"github.com/kun1ts4/link-shortener/internal/domain"
)

// Handler handles HTTP requests for the link-shortener API.
type Handler struct {
	uc  domain.UseCase
	log *slog.Logger
	cfg *config.RateLimitConfig
}

// NewHandler creates a new Handler instance.
func NewHandler(uc domain.UseCase, log *slog.Logger, cfg *config.RateLimitConfig) *Handler {
	return &Handler{
		uc:  uc,
		log: log,
		cfg: cfg,
	}
}

// errorResponse represents an error response body.
type errorResponse struct {
	Error string `json:"error"`
}

// writeError writes an error response with the given status code and message.
func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	b, _ := json.Marshal(errorResponse{Error: msg})
	_, _ = w.Write(b)
}
