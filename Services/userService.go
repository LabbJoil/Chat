package service

import (
	"log"

	remotemodels "github.com/LabbJoil/Chat/Models/RemoteModels"
)

func (SC *ChatInteraction) GetUser(idUser string) (*remotemodels.UserInfo, error) {
	user, err := SC.DBConnection.GetUser(idUser)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("Get user with id: %v\n", user.ID)
	return &remotemodels.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (SC *ChatInteraction) CreateUser(userName, email string) (*remotemodels.UserInfo, error) {
	user, err := SC.DBConnection.InsertNewUser(userName, email)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	log.Printf("User with ID %s added\n", user.ID)
	return &remotemodels.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (SC *ChatInteraction) ConcateConversation(idUser, idChat string) error {
	err := SC.DBConnection.ConcatConversation(idUser, idChat)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("User with ID %s concat the chat with id %s\n", idUser, idChat)
	return nil
}

func (SC *ChatInteraction) GetConversations(idUser string) ([]*remotemodels.Chat, error) {
	allUserChatsDBM, err := SC.DBConnection.GetAllConversations(idUser)
	if err != nil {
		return nil, err
	}
	log.Printf("Get %d chats from a user with id %s", len(allUserChatsDBM), idUser)
	var chats []*remotemodels.Chat
	for key, value := range allUserChatsDBM {
		chat := remotemodels.Chat{
			Title:        key.Title,
			CountMembers: value,
		}
		chats = append(chats, &chat)
	}
	return chats, nil
}

func (SC *ChatInteraction) LeaveChat(idUser, idConversation string) error {
	err := SC.DBConnection.LeaveConversation(idUser, idConversation)
	if err != nil {
		return err
	}
	log.Printf("User with id: %s leave conversation with id: %s\n", idUser, idConversation)
	return nil
}

func (SC *ChatInteraction) DeleteUser(idUser string) error {
	err := SC.DBConnection.DeleteUser(idUser)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("User with ID %s delete\n", idUser)
	return nil
}
