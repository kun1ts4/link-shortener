package memory

import (
	"context"
	"link-shortener/internal/domain"
	"sync"
)

type MemRepository struct {
	m  sync.RWMutex
	db map[string]domain.Link
}

func NewMemRepository() *MemRepository {
	return &MemRepository{
		m:  sync.RWMutex{},
		db: make(map[string]domain.Link),
	}
}

func (r *MemRepository) Create(_ context.Context, link *domain.Link) error {
	r.m.Lock()
	defer r.m.Unlock()
	_, ok := r.db[link.Short]
	if ok {
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

	result := &domain.Link{
		Original:  link.Original,
		Short:     link.Short,
		Clicks:    link.Clicks,
		CreatedAt: link.CreatedAt,
	}

	return result, nil
}

func (r *MemRepository) FindByOriginal(_ context.Context, original string) (*domain.Link, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	for _, link := range r.db {
		if link.Original == original {
			result := &domain.Link{
				Original:  link.Original,
				Short:     link.Short,
				Clicks:    link.Clicks,
				CreatedAt: link.CreatedAt,
			}
			return result, nil
		}
	}

	return nil, domain.ErrNotFound
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
