package models

import "time"

type Purchase struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Item      string
	Price     uint
	CreatedAt time.Time
}
