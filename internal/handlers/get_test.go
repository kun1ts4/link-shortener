package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kun1ts4/link-shortener/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLink_Redirect(t *testing.T) {
	link := &domain.Link{
		Original:  "https://google.com",
		Short:     "abcde12345",
		CreatedAt: time.Now(),
	}

	uc := &mockUseCase{
		getOriginalLinkFunc: func(_ context.Context, short string) (*domain.Link, error) {
			assert.Equal(t, "abcde12345", short)
			return link, nil
		},
	}

	h := newTestHandler(uc)
	srv := httptest.NewServer(h.NewRouter())
	defer srv.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(srv.URL + "/abcde12345")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode)
	assert.Equal(t, link.Original, resp.Header.Get("Location"))
}

func TestGetLink_JSON(t *testing.T) {
	link := &domain.Link{
		Original:  "https://google.com",
		Short:     "abcde12345",
		CreatedAt: time.Now(),
	}

	uc := &mockUseCase{
		getOriginalLinkFunc: func(_ context.Context, short string) (*domain.Link, error) {
			assert.Equal(t, "abcde12345", short)
			return link, nil
		},
	}

	h := newTestHandler(uc)
	srv := httptest.NewServer(h.NewRouter())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/abcde12345.json")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var body GetLinkJsonResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, link.Original, body.URL)
}

func TestGetLink_NotFound(t *testing.T) {
	uc := &mockUseCase{
		getOriginalLinkFunc: func(_ context.Context, _ string) (*domain.Link, error) {
			return nil, domain.ErrNotFound
		},
	}

	h := newTestHandler(uc)
	srv := httptest.NewServer(h.NewRouter())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/notexist")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
