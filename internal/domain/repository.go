package domain

type Repository interface {
	Create(link *Link) error
	FindByShort(key string) (*Link, error)
	FindByID(id int64) (*Link, error)
	FindByOriginal(original string) (*Link, error)
	IncrementClicks(id int64) error
}
