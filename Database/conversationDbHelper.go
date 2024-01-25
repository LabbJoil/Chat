package databasePostgreSQL

import (
	"fmt"
	"log"
	"time"

	databasemodels "github.com/LabbJoil/Chat/Models/DatabaseModels"
	"github.com/google/uuid"
)

func (DBD *DatabaseDevelop) CreateConversation(name, idCreator string) (*databasemodels.Conversation, error) {
	newChat := &databasemodels.Conversation{
		ID:        "chat_" + uuid.New().String(),
		Title:     name,
		IDCreator: idCreator,
		CreatedAt: time.Now(),
	}
	result := DBD.DB.Create(&newChat)
	if result.Error != nil {
		return nil, result.Error
	}
	return newChat, nil
}

func (DBD *DatabaseDevelop) ConcatConversation(idUser, idConversation string) error {
	var chat databasemodels.Conversation
	var user databasemodels.User
	var countConcat int64

	if result := DBD.DB.Where("id = ?", idConversation).First(&chat); result.Error != nil {
		return result.Error
	}
	if result := DBD.DB.Where("id = ?", idUser).First(&user); result.Error != nil {
		return result.Error
	}
	if result := DBD.DB.Table("conversation_users").Where("id_conversation = ? AND id_user = ?", idConversation, idUser).Count(&countConcat); result.Error != nil {
		return result.Error
	} else if countConcat != 0 {
		newErr := fmt.Errorf("Any error. User with id: %s has already connected to the conversation with id: %s", idUser, idConversation)
		return newErr
	}
	chatUser := databasemodels.ConversationUser{
		IDUser:         user.ID,
		IDConversation: chat.ID,
	}
	DBD.DB.Create(&chatUser)
	return nil
}

func (DBD *DatabaseDevelop) LeaveConversation(idUser, idConversation string) error {
	var conversationUser databasemodels.ConversationUser
	result := DBD.DB.Where("id_user = ? AND id_conversation = ?", idUser, idConversation).First(&conversationUser)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}

	result = DBD.DB.Delete(&conversationUser)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		newErr := fmt.Errorf("Any error. User with id: %s not leave conversation with id: %s", idUser, idConversation)
		return newErr
	}
	return nil
}

func (DBD *DatabaseDevelop) DeleteConversation(idCreator, idConversation string) error {
	var conversation databasemodels.Conversation
	result := DBD.DB.Where("id = ? AND id_creator = ?", idConversation, idCreator).First(&conversation)
	if result.Error != nil {
		return result.Error
	}
	result = DBD.DB.Delete(&conversation)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		newErr := fmt.Errorf("Any error. Conversation with id: %s not deleted by creator with id: %s", idConversation, idCreator)
		return newErr
	}
	return nil
}

func (DBD *DatabaseDevelop) GetAllConversations(idUser string) (map[databasemodels.Conversation]int, error) {
	var conversationIDs []string
	result := DBD.DB.Table("conversation_users").Where("id_user = ?", idUser).Pluck("id_conversation", &conversationIDs)
	if result.Error != nil {
		return nil, result.Error
	}

	conversationMembersDict := make(map[databasemodels.Conversation]int)
	for _, idConversation := range conversationIDs {
		var countUsers int64
		var conversation databasemodels.Conversation
		result := DBD.DB.Table("conversation_users").Where("id_conversation = ?", idConversation).Count(&countUsers)
		if result.Error != nil {
			return nil, result.Error
		}
		result = DBD.DB.Table("conversations").Where("Id = ?", idConversation).Scan(&conversation)
		if result.Error != nil {
			return nil, result.Error
		}
		conversationMembersDict[conversation] = int(countUsers)
	}
	return conversationMembersDict, nil
}

func (DBD *DatabaseDevelop) GetConversationMembers(idConversation string) ([]*databasemodels.User, error) {
	var membersId []string
	result := DBD.DB.Table("conversation_users").
		Where("id_conversation = ?", idConversation).
		Pluck("id_user", &membersId)
	if result.Error != nil {
		return nil, result.Error
	}
	var membersChat []*databasemodels.User
	for _, userId := range membersId {
		var user *databasemodels.User
		result := DBD.DB.Table("users").Where("Id = ?", userId).Scan(&user)
		if result.Error != nil {
			return nil, result.Error
		}
		membersChat = append(membersChat, user)
	}
	return membersChat, nil
}
