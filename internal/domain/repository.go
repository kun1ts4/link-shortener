package domain

type Repository interface {
	Create(link *Link) error
	FindByShort(short string) (*Link, error)
	FindByOriginal(original string) (*Link, error)
	IncrementClicks(short string) error
}
