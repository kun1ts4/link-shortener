package memory

import (
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

func (r *MemRepository) Create(link *domain.Link) error {
	r.m.Lock()
	defer r.m.Unlock()
	_, ok := r.db[link.Short]
	if ok == true {
		return domain.ErrAlreadyExists
	}

	r.db[link.Short] = *link

	return nil
}

func (r *MemRepository) FindByShort(short string) (*domain.Link, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	link, ok := r.db[short]
	if !ok {
		return &domain.Link{}, domain.ErrNotFound
	}

	result := &domain.Link{
		Original:  link.Original,
		Short:     link.Short,
		Clicks:    link.Clicks,
		CreatedAt: link.CreatedAt,
	}

	return result, nil
}

func (r *MemRepository) FindByOriginal(original string) (*domain.Link, error) {
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

	return &domain.Link{}, domain.ErrNotFound
}

func (r *MemRepository) IncrementClicks(short string) error {
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
