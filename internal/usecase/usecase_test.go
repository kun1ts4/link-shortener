package usecase

import (
	"context"
	"link-shortener/internal/domain"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRepo struct {
	db map[string]*domain.Link
}

func newMockRepo() *mockRepo {
	return &mockRepo{db: make(map[string]*domain.Link)}
}

func (m *mockRepo) Create(_ context.Context, link *domain.Link) error {
	if _, ok := m.db[link.Short]; ok {
		return domain.ErrAlreadyExists
	}
	m.db[link.Short] = link
	return nil
}

func (m *mockRepo) FindByShort(_ context.Context, short string) (*domain.Link, error) {
	l, ok := m.db[short]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return l, nil
}

func (m *mockRepo) IncrementClicks(_ context.Context, short string) error {
	l, ok := m.db[short]
	if !ok {
		return domain.ErrNotFound
	}
	l.Clicks++
	return nil
}

type mockGen struct {
	results []string
	attempt int
}

func (g *mockGen) Generate(_ string) string {
	s := g.results[g.attempt]
	if g.attempt < len(g.results)-1 {
		g.attempt++
	}
	return s
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func TestCreateShort_NewLink(t *testing.T) {
	repo := newMockRepo()
	gen := &mockGen{results: []string{"abcde12345"}}
	uc := NewShortenerUseCase(gen, repo, newLogger())

	link, err := uc.CreateShort(context.Background(), "https://google.com")
	require.NoError(t, err)
	assert.Equal(t, "https://google.com", link.Original)
	assert.Equal(t, "abcde12345", link.Short)
}

func TestCreateShort_Collision(t *testing.T) {
	repo := newMockRepo()

	existing := &domain.Link{Original: "https://ya.ru", Short: "short1", CreatedAt: time.Now()}
	repo.db["short1"] = existing

	gen := &mockGen{results: []string{"short1", "short2"}}
	uc := NewShortenerUseCase(gen, repo, newLogger())

	link, err := uc.CreateShort(context.Background(), "https://google.com")
	require.NoError(t, err)
	assert.Equal(t, "short2", link.Short)
	assert.Len(t, repo.db, 2)
}
