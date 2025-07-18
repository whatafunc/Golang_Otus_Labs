package storage

import "time"

type Event struct {
	ID          int // auto-increment or assigned
	Title       string
	Description string
	Start       *time.Time // nullable
	End         *time.Time // nullable
	AllDay      float64
	Clinic      *string // nullable
	UserID      *int    // nullable
	Service     *string // nullable
}
