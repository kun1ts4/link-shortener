package postgres

import (
	"time"

	"github.com/kun1ts4/link-shortener/internal/domain"
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
		Original:  l.Original,
		Short:     l.Short,
		Clicks:    l.Clicks,
		CreatedAt: l.CreatedAt,
	}
}

func FromDomain(link *domain.Link) *LinkDB {
	return &LinkDB{
		Original:  link.Original,
		Short:     link.Short,
		Clicks:    link.Clicks,
		CreatedAt: link.CreatedAt,
	}
}
