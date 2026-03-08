package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"link-shortener/internal/domain"
	"net/http"
)

type CreateLinkRequest struct {
	URL string `json:"url"`
}

type CreateLinkResponse struct {
	Short    string `json:"short"`
	Original string `json:"original"`
}

func (h *Handler) CreateLink(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Info("reading body", "error", err)
		writeError(w, http.StatusBadRequest, "error read request")
		return
	}

	req := CreateLinkRequest{}
	if err = json.Unmarshal(body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	link, err := h.uc.CreateShort(r.Context(), req.URL)
	if err != nil {
		if errors.Is(err, domain.ErrInvalid) {
			writeError(w, http.StatusBadRequest, "invalid URL")
			return
		} else if errors.Is(err, domain.ErrStorageFull) {
			writeError(w, http.StatusServiceUnavailable, "storage full")
			return
		}
		h.log.Info("create short failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp, err := json.Marshal(CreateLinkResponse{
		Short:    link.Short,
		Original: link.Original,
	})
	if err != nil {
		h.log.Info("marshalling response", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		h.log.Info("failed to write response", "error", err)
	}

}
