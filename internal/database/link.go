package database

import (
	"time"
)

// Link is the persisted short link.
type Link struct {
	ID        uint      `gorm:"primaryKey"`
	Code      string    `gorm:"size:7;uniqueIndex;not null"`
	LongURL   string    `gorm:"size:768;uniqueIndex;not null"` // 768*4 bytes = 3072, utf8mb4 index limit
	CreatedAt time.Time
}
