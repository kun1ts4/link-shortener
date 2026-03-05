package domain

import (
	"time"
)

type Link struct {
	ID        int64
	Original  string
	Short     string
	Clicks    int64
	CreatedAt time.Time
}
