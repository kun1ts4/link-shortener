package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/kun1ts4/link-shortener/internal/domain"
)

type CreateLinkRequest struct {
	URL string `json:"url"`
}

type CreateLinkResponse struct {
	Short    string `json:"short"`
	Original string `json:"original"`
}

func (h *Handler) CreateLink(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("request-id")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error("failed to read request body",
			"request_id", requestID,
			"error", err)
		writeError(w, http.StatusBadRequest, "error read request")
		return
	}

	req := CreateLinkRequest{}
	if err = json.Unmarshal(body, &req); err != nil {
		h.log.Warn("invalid JSON in request body",
			"request_id", requestID)
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// Validate URL field presence
	if req.URL == "" {
		h.log.Warn("missing url field in request",
			"request_id", requestID)
		writeError(w, http.StatusBadRequest, "url field is required")
		return
	}

	link, err := h.uc.CreateShort(r.Context(), req.URL)
	if err != nil {
		if errors.Is(err, domain.ErrInvalid) {
			h.log.Warn("invalid URL provided",
				"request_id", requestID,
				"url", req.URL)
			writeError(w, http.StatusBadRequest, "invalid URL")
			return
		} else if errors.Is(err, domain.ErrStorageFull) {
			h.log.Error("storage is full",
				"request_id", requestID)
			writeError(w, http.StatusServiceUnavailable, "storage full")
			return
		}
		h.log.Error("failed to create short link",
			"request_id", requestID,
			"error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.log.Info("short link created successfully",
		"request_id", requestID,
		"short", link.Short)

	resp, err := json.Marshal(CreateLinkResponse{
		Short:    link.Short,
		Original: link.Original,
	})
	if err != nil {
		h.log.Error("failed to marshal response",
			"request_id", requestID,
			"error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		h.log.Error("failed to write response",
			"request_id", requestID,
			"error", err)
	}

}
