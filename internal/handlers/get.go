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
	URL string `json:"url"`
}

func (h *Handler) GetLink(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "short")
	format, _ := r.Context().Value(middleware.URLFormatCtxKey).(string)

	link, err := h.uc.GetOriginalLink(r.Context(), short)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			h.log.Info("not found", "error", err)
			w.WriteHeader(http.StatusNotFound)
			return
		} else if errors.Is(err, domain.ErrInvalid) {
			h.log.Info("invalid", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.log.Info("error getting link", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if format == "json" {
		resp, err := json.Marshal(&GetLinkJsonResponse{
			URL: link.Original,
		})
		if err != nil {
			h.log.Error("error marshaling link", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
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
