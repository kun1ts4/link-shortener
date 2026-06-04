package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/kun1ts4/link-shortener/internal/config"
	"github.com/kun1ts4/link-shortener/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUseCase struct {
	createShortFunc     func(ctx context.Context, original string) (*domain.Link, error)
	getOriginalLinkFunc func(ctx context.Context, short string) (*domain.Link, error)
}

func (m *mockUseCase) CreateShort(ctx context.Context, original string) (*domain.Link, error) {
	return m.createShortFunc(ctx, original)
}

func (m *mockUseCase) GetOriginalLink(ctx context.Context, short string) (*domain.Link, error) {
	return m.getOriginalLinkFunc(ctx, short)
}

func newTestHandler(uc domain.UseCase) *Handler {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &config.RateLimitConfig{RequestsPerSecond: 100}
	return NewHandler(uc, log, cfg)
}

func TestCreateLink_OK(t *testing.T) {
	link := &domain.Link{
		Original:  "https://google.com",
		Short:     "abcde12345",
		CreatedAt: time.Now(),
	}

	uc := &mockUseCase{
		createShortFunc: func(_ context.Context, original string) (*domain.Link, error) {
			return link, nil
		},
	}

	h := newTestHandler(uc)
	body, _ := json.Marshal(CreateLinkRequest{URL: "https://google.com"})
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateLink(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp CreateLinkResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, link.Short, resp.Short)
	assert.Equal(t, link.Original, resp.Original)
}

func TestCreateLink_InvalidJSON(t *testing.T) {
	h := newTestHandler(&mockUseCase{})

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("not json"))
	w := httptest.NewRecorder()

	h.CreateLink(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateLink_InvalidURL(t *testing.T) {
	uc := &mockUseCase{
		createShortFunc: func(_ context.Context, _ string) (*domain.Link, error) {
			return nil, domain.ErrInvalid
		},
	}

	h := newTestHandler(uc)
	body, _ := json.Marshal(CreateLinkRequest{URL: "noturl"})
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateLink(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateLink_StorageFull(t *testing.T) {
	uc := &mockUseCase{
		createShortFunc: func(_ context.Context, _ string) (*domain.Link, error) {
			return nil, domain.ErrStorageFull
		},
	}

	h := newTestHandler(uc)
	body, _ := json.Marshal(CreateLinkRequest{URL: "https://google.com"})
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateLink(w, r)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}
