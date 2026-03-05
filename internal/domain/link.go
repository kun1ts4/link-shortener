package domain

import (
	"time"
)

type Link struct {
	Original  string
	Short     string
	Clicks    int64
	CreatedAt time.Time
}
