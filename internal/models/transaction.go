package models

import "time"

type Transaction struct {
	ID         uint `gorm:"primaryKey"`
	SenderID   uint
	ReceiverID uint
	Amount     int
	CreatedAt  time.Time
}
