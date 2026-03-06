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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := CreateLinkRequest{}
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link, err := h.uc.CreateShort(r.Context(), req.URL)
	if err != nil {
		if errors.Is(err, domain.ErrInvalid) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.log.Info("create short failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(CreateLinkResponse{
		Short:    link.Short,
		Original: link.Original,
	})
	if err != nil {
		h.log.Info("marshalling response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		h.log.Info("failed to write response", "error", err)
	}

}
