package domain

import "context"

type Repository interface {
	Create(ctx context.Context, link *Link) error
	FindByShort(ctx context.Context, short string) (*Link, error)
	IncrementClicks(ctx context.Context, short string) error
}
