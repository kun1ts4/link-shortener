package domain

import "context"

type UseCase interface {
	CreateShort(ctx context.Context, original string) (*Link, error)
	GetOriginalLink(ctx context.Context, short string) (*Link, error)
}
