package postgres

import (
	"link-shortener/internal/domain"
	"time"
)

type LinkDB struct {
	ID        int64     `db:"id"`
	Original  string    `db:"original"`
	Short     string    `db:"short"`
	Clicks    int64     `db:"clicks"`
	CreatedAt time.Time `db:"created_at"`
}

func (l *LinkDB) ToDomain() *domain.Link {
	return &domain.Link{
		ID:        l.ID,
		Original:  l.Original,
		Short:     l.Short,
		Clicks:    l.Clicks,
		CreatedAt: l.CreatedAt,
	}
}

func FromDomain(link *domain.Link) *LinkDB {
	return &LinkDB{
		ID:        link.ID,
		Original:  link.Original,
		Short:     link.Short,
		Clicks:    link.Clicks,
		CreatedAt: link.CreatedAt,
	}
}
