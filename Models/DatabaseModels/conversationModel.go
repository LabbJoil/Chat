package databasemodels

import "time"

type Conversation struct {
	ID        string    `gorm:"primarykey"`
	Title     string    `gorm:"not null"`
	IDCreator string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
