package handlers

import (
	"link-shortener/internal/config"
	"link-shortener/internal/domain"
	"log/slog"
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
