package databasePostgreSQL

import (
	"fmt"
	"log"

	databasemodels "github.com/LabbJoil/Chat/Models/DatabaseModels"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func (DBD *DatabaseDevelop) InsertNewMessage(idSender, idConversation, content string) (*databasemodels.Message, error) {
	message := &databasemodels.Message{
		ID:             "message_" + uuid.New().String(),
		IDSender:       idSender,
		Content:        content,
		IDConversation: idConversation,
	}
	result := DBD.DB.Create(message)
	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}
	return message, nil
}

func (DBD *DatabaseDevelop) DeleteMessage(idMessage, idUser string) error {
	var message databasemodels.Message
	result := DBD.DB.Where("id = ? AND id_sender = ?", idMessage, idUser).First(&message)
	if result.Error != nil {
		return result.Error
	}
	result = DBD.DB.Delete(&message)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		newErr := fmt.Errorf("Any error. Message with id: %s not delete", idMessage)
		return newErr
	}
	return nil
}

func (DBD *DatabaseDevelop) GetAllMessages(idUser, idConversation string) ([]*databasemodels.Message, error) {
	var messages []*databasemodels.Message
	query := DBD.DB.
		Table("messages").
		Select("messages.id, content, id_sender, conversations.id, messages.created_at").
		Joins("JOIN conversations ON messages.id_conversation = conversations.id").
		Joins("JOIN conversation_users ON conversations.id = conversation_users.id_conversation").
		Where("conversations.id = ? AND conversation_users.id_user = ?", idConversation, idUser).
		Order("messages.created_at").
		Scan(&messages)

	if query.Error != nil {
		log.Println(query.Error)
		return nil, query.Error
	}
	return messages, nil
}

func (DBD *DatabaseDevelop) GetMessages(idUser string, conversationIds, messagesIds []string) (*[]databasemodels.Message, error) {
	var messages []databasemodels.Message

	query := DBD.DB.
		Table("messages").
		Select("messages.id, messages.content, messages.id_sender, messages.created_at, messages.id_conversation").
		Joins("JOIN conversations ON messages.id_conversation = conversations.id").
		Joins("JOIN conversation_users ON conversation_users.id_conversation = conversations.id").
		Where("conversation_users.id_user = ? AND conversation_users.id_conversation IN ?", idUser, conversationIds)

	if len(messagesIds) > 0 {
		query = query.Not("messages.id IN ?", messagesIds)
	}
	if err := query.Scan(&messages).Error; err != nil {
		return nil, err
	}
	return &messages, nil
}

func (DBD *DatabaseDevelop) GetFilteredMessages(idUser, idConversation, idMessages string, countMessages int) ([]*databasemodels.Message, error) {
	var messages []*databasemodels.Message
	query := DBD.DB.
		Table("messages").
		Select("messages.id, messages.content, messages.id_sender, messages.id_conversation, messages.created_at").
		Joins("JOIN conversations ON messages.id_conversation = conversations.id").
		Joins("JOIN conversation_users ON conversations.id = conversation_users.id_conversation").
		Where("conversations.id = ? AND conversation_users.id_user = ? AND messages.created_at >= (SELECT created_at FROM messages WHERE id = ?)", idConversation, idUser, idMessages).
		Order("messages.created_at").
		Limit(countMessages).
		Scan(&messages)

	if query.Error != nil {
		return nil, query.Error
	}
	return messages, nil
}
