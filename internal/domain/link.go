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
	u, err := url.ParseRequestURI(original)
	if err != nil {
		return nil, ErrInvalid
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, ErrInvalid
	}
	if u.Host == "" {
		return nil, ErrInvalid
	}

	return &Link{
		Original:  original,
		Short:     short,
		Clicks:    0,
		CreatedAt: time.Now(),
	}, nil
}
