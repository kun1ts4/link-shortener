package handlers

import (
	"link-shortener/internal/domain"
	"log/slog"
)

type Handler struct {
	uc  domain.UseCase
	log *slog.Logger
}

func NewHandler(uc domain.UseCase, log *slog.Logger) *Handler {
	return &Handler{
		uc:  uc,
		log: log,
	}
}
