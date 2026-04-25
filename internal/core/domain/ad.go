package domain

import "time"

type Ad struct {
	ID          int
	Version     int
	Title       string
	Description string
	Price       float64
	UserID      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
