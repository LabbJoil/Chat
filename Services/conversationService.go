package service

import (
	"fmt"
	"log"

	remotemodels "github.com/LabbJoil/Chat/Models/RemoteModels"
)

func (SC *ChatInteraction) CreateConversation(title, userId string) (*remotemodels.Chat, error) {
	conversation, err := SC.DBConnection.CreateConversation(title, userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("Chat with ID %s create\n", conversation.ID)
	err = SC.DBConnection.ConcatConversation(userId, conversation.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("User with ID %s concat the chat with id %s\n", userId, conversation.ID)
	return &remotemodels.Chat{
		ID:           conversation.ID,
		Title:        conversation.Title,
		CountMembers: 1,
	}, nil
}

func (SC *ChatInteraction) DeleteConversation(idUser, idConversation string) error {
	if err := SC.DBConnection.DeleteConversation(idUser, idConversation); err != nil {
		log.Print(err)
		return err
	}
	log.Printf("Conversation with id: %s deleted by creator with id: %s\n", idConversation, idUser)
	return nil
}

func (SC *ChatInteraction) GetConversationMembers(idUser, idConversation string) ([]*remotemodels.UserInfo, error) {
	members, err := SC.DBConnection.GetConversationMembers(idConversation)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("Get %d members from chat with id %s", len(members), idConversation)

	if !IsUserInSlice(idUser, members) {
		newErr := fmt.Errorf("Any error. User: %s not in chat", idUser)
		log.Println(newErr)
		return nil, newErr
	}

	var membersChat []*remotemodels.UserInfo
	for _, userDb := range members {
		user := remotemodels.UserInfo{
			ID:       userDb.ID,
			Username: userDb.Username,
			Email:    userDb.Email,
		}
		membersChat = append(membersChat, &user)
	}
	return membersChat, nil
}
