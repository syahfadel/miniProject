package entities

import "time"

//
type Post struct {
	ID        uint `gorm:"primaryKey"`
	Body      string
	UserID    uint
	CreatedAt time.Time
}
