package handlers

import (
	"encoding/json"
	"errors"
	"link-shortener/internal/domain"
	"net/http"

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

	http.Redirect(w, r, link.Original, http.StatusMovedPermanently)
}
