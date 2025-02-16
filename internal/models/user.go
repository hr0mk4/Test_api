package models

type User struct {
	ID                   uint   `gorm:"primaryKey"`
	Username             string `gorm:"uniqueIndex"`
	Password             string
	Coins                uint
	Purchases            []Purchase
	SentTransactions     []Transaction `gorm:"foreignKey:SenderID"`
	ReceivedTransactions []Transaction `gorm:"foreignKey:ReceiverID"`
}
