package databasemodels

type ConversationUser struct {
	ID             int `gorm:"primaryKey"`
	IDConversation string
	Conv           Conversation `gorm:"foreignKey:IDConversation;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"conversation"`
	IDUser         string
	Us             User `gorm:"foreignKey:IDUser;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
}
