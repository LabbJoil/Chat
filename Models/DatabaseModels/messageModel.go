package databasemodels

import "time"

type Message struct {
	ID             string       `gorm:"primarykey"`
	Content        string       `gorm:"not null"`
	IDSender       string       `gorm:"not null"`
	CreatedAt      time.Time    `gorm:"default:CURRENT_TIMESTAMP"`
	IDConversation string       `gorm:"many2many:conversation"`
	Conversation   Conversation `gorm:"foreignKey:IDConversation;;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"conversation"`
}
