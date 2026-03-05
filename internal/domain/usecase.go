package domain

type UseCase interface {
	CreateShort(original string) (*Link, error)
	GetOriginalUrl(key string) (*Link, error)
}
