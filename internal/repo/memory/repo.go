package memory

import (
	"context"
	"sync"

	"github.com/kun1ts4/link-shortener/internal/config"
	"github.com/kun1ts4/link-shortener/internal/domain"
)

type MemRepository struct {
	m       sync.RWMutex
	db      map[string]domain.Link
	maxSize int
}

func NewMemRepository(cfg config.MemoryConfig) *MemRepository {
	return &MemRepository{
		db:      make(map[string]domain.Link),
		maxSize: cfg.MaxSize,
	}
}

func (r *MemRepository) Create(_ context.Context, link *domain.Link) error {
	r.m.Lock()
	defer r.m.Unlock()

	if r.maxSize > 0 && len(r.db) >= r.maxSize {
		return domain.ErrStorageFull
	}

	if _, ok := r.db[link.Short]; ok {
		return domain.ErrAlreadyExists
	}

	r.db[link.Short] = *link

	return nil
}

func (r *MemRepository) FindByShort(_ context.Context, short string) (*domain.Link, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	link, ok := r.db[short]
	if !ok {
		return nil, domain.ErrNotFound
	}

	return &link, nil
}

func (r *MemRepository) IncrementClicks(_ context.Context, short string) error {
	r.m.Lock()
	defer r.m.Unlock()
	link, ok := r.db[short]
	if !ok {
		return domain.ErrNotFound
	}
	link.Clicks++
	r.db[short] = link

	return nil
}
