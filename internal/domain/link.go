package domain

import (
	"net/url"
	"time"
)

type Link struct {
	Original  string
	Short     string
	Clicks    int64
	CreatedAt time.Time
}

func NewLink(original, short string) (*Link, error) {
	_, err := url.Parse(original)
	if err != nil {
		return nil, ErrInvalid
	}

	return &Link{
		Original:  original,
		Short:     short,
		Clicks:    0,
		CreatedAt: time.Now(),
	}, nil
}
