package memory

import (
	"context"
	"link-shortener/internal/config"
	"link-shortener/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryRepoCreate(t *testing.T) {
	ctx := context.Background()
	repo := NewMemRepository(config.MemoryConfig{})

	link := &domain.Link{
		Original:  "https://google.com",
		Short:     "abc012",
		Clicks:    0,
		CreatedAt: time.Now(),
	}

	err := repo.Create(ctx, link)
	require.NoError(t, err)

	found, err := repo.FindByShort(ctx, "abc012")
	require.NoError(t, err)

	assert.Equal(t, link.Original, found.Original)
	assert.Equal(t, link.Short, found.Short)
	assert.Equal(t, link.Clicks, found.Clicks)
}

func TestMemoryRepoCreateDuplicate(t *testing.T) {
	ctx := context.Background()
	repo := NewMemRepository(config.MemoryConfig{})

	link := &domain.Link{
		Original:  "https://google.com",
		Short:     "abc012",
		Clicks:    0,
		CreatedAt: time.Now(),
	}

	err := repo.Create(ctx, link)
	require.NoError(t, err)

	err = repo.Create(ctx, link)
	assert.ErrorIs(t, err, domain.ErrAlreadyExists)
}

func TestMemoryRepoStorageFull(t *testing.T) {
	ctx := context.Background()
	repo := NewMemRepository(config.MemoryConfig{MaxSize: 1})

	err := repo.Create(ctx, &domain.Link{Original: "https://google.com", Short: "short00001", CreatedAt: time.Now()})
	require.NoError(t, err)

	err = repo.Create(ctx, &domain.Link{Original: "https://yandex.ru", Short: "short00002", CreatedAt: time.Now()})
	assert.ErrorIs(t, err, domain.ErrStorageFull)
}

func TestMemoryRepoFindByOriginal(t *testing.T) {
	ctx := context.Background()
	repo := NewMemRepository(config.MemoryConfig{})

	link := &domain.Link{
		Original:  "https://google.com",
		Short:     "abc012",
		Clicks:    0,
		CreatedAt: time.Now(),
	}

	err := repo.Create(ctx, link)
	require.NoError(t, err)

	found, err := repo.FindByOriginal(ctx, "https://google.com")
	require.NoError(t, err)

	assert.Equal(t, link.Original, found.Original)
	assert.Equal(t, link.Short, found.Short)
	assert.Equal(t, link.Clicks, found.Clicks)
}

func TestMemoryRepoIncrement(t *testing.T) {
	ctx := context.Background()
	repo := NewMemRepository(config.MemoryConfig{})

	link := &domain.Link{
		Original:  "https://google.com",
		Short:     "abc012",
		Clicks:    0,
		CreatedAt: time.Now(),
	}

	err := repo.Create(ctx, link)
	require.NoError(t, err)

	err = repo.IncrementClicks(ctx, "abc012")
	require.NoError(t, err)

	found, err := repo.FindByShort(ctx, "abc012")
	require.NoError(t, err)

	assert.Equal(t, link.Original, found.Original)
	assert.Equal(t, link.Short, found.Short)
	assert.Equal(t, link.Clicks+1, found.Clicks)
}

func TestMemoryRepoNotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewMemRepository(config.MemoryConfig{})

	_, err := repo.FindByShort(ctx, "not exist")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = repo.FindByOriginal(ctx, "https://not.exist")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	err = repo.IncrementClicks(ctx, "not exist")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
