package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kun1ts4/link-shortener/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type GetLinkJsonResponse struct {
	URL    string `json:"url"`
	Clicks int64  `json:"clicks"`
}

func (h *Handler) GetLink(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "short")
	format, _ := r.Context().Value(middleware.URLFormatCtxKey).(string)

	if err := validateURL(short); err != nil {
		h.log.Info("invalid short link", "short", short, "error", err)
		writeError(w, http.StatusBadRequest, "invalid short link")
		return
	}

	link, err := h.uc.GetOriginalLink(r.Context(), short)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			h.log.Info("not found", "error", err)
			writeError(w, http.StatusNotFound, "link not found")
			return
		} else if errors.Is(err, domain.ErrInvalid) {
			h.log.Info("invalid", "error", err)
			writeError(w, http.StatusBadRequest, "invalid short link")
			return
		}

		h.log.Info("error getting link", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if format == "json" {
		resp, err := json.Marshal(&GetLinkJsonResponse{
			URL:    link.Original,
			Clicks: link.Clicks,
		})
		if err != nil {
			h.log.Error("error marshaling link", "error", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resp)
		if err != nil {
			h.log.Error("failed to write response", "error", err)
		}
		return
	}

	writeError(w, http.StatusBadRequest, "use .json format to get link")
}

func validateURL(url string) error {
	if len(url) != 10 {
		return fmt.Errorf("invalid url")
	}

	return nil
}
