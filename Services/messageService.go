package service

import (
	"fmt"
	"log"

	remotemodels "github.com/LabbJoil/Chat/Models/RemoteModels"
)

func (SC *ChatInteraction) CreateMessage(idSender, idConversation, content string) (*remotemodels.ChatMessage, error) {
	membersChat, err := SC.DBConnection.GetConversationMembers(idConversation)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("Get members with ids: %v\n", membersChat)
	if !IsUserInSlice(idSender, membersChat) {
		newErr := fmt.Errorf("Any error. User: %s not in chat", idSender)
		log.Println(newErr)
		return nil, newErr
	}
	newMessage, err := SC.DBConnection.InsertNewMessage(idSender, idConversation, content)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("Message with id: %s create \n", newMessage.ID)
	return &remotemodels.ChatMessage{
		ID:        newMessage.ID,
		Content:   newMessage.Content,
		IDSender:  newMessage.IDSender,
		IDChat:    newMessage.IDConversation,
		CreatedAt: newMessage.CreatedAt.String(),
	}, nil
}

func (SC *ChatInteraction) DeleteMessage(idMessage, idUser string) error {
	if err := SC.DBConnection.DeleteMessage(idMessage, idUser); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Message with id: %s delete \n", idMessage)
	return nil
}

func (SC *ChatInteraction) GetAllMessages(userId, conversationID string) ([]*remotemodels.ChatMessage, error) {
	messagesDBM, err := SC.DBConnection.GetAllMessages(userId, conversationID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var messages []*remotemodels.ChatMessage
	for _, messageDBM := range messagesDBM {
		messsageRM := remotemodels.ChatMessage{
			ID:        messageDBM.ID,
			Content:   messageDBM.Content,
			IDSender:  messageDBM.IDSender,
			CreatedAt: messageDBM.CreatedAt.String(),
			IDChat:    conversationID,
		}
		messages = append(messages, &messsageRM)
	}
	log.Printf("Get: %d all messages \n", len(messages))
	return messages, nil
}

func (SC *ChatInteraction) GetFilterMessages(userId, conversationID, messageId string, count int) ([]*remotemodels.ChatMessage, error) {
	messagesDBM, err := SC.DBConnection.GetFilteredMessages(userId, conversationID, messageId, count)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var messages []*remotemodels.ChatMessage
	for _, messageDBM := range messagesDBM {
		messsageRM := remotemodels.ChatMessage{
			ID:        messageDBM.ID,
			Content:   messageDBM.Content,
			IDSender:  messageDBM.IDSender,
			CreatedAt: messageDBM.CreatedAt.String(),
			IDChat:    conversationID,
		}
		messages = append(messages, &messsageRM)
	}
	log.Printf("Get: %d filtered messages \n", len(messages))
	return messages, nil
}

func (SC *ChatInteraction) ThreadGetMessages(idUser string) (*[]remotemodels.ChatMessage, error) {
	allUserChatsDBM, err := SC.DBConnection.GetAllConversations(idUser)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var chatIds []string
	for key := range allUserChatsDBM {
		chatIds = append(chatIds, key.ID)
	}
	messages, err := SC.DBConnection.GetMessages(idUser, chatIds, SC.ReceivedMessages[idUser])
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var newMessages []remotemodels.ChatMessage
	for _, message := range *messages {
		newMessage := remotemodels.ChatMessage{
			ID:        message.ID,
			Content:   message.Content,
			IDSender:  message.IDSender,
			IDChat:    message.IDConversation,
			CreatedAt: message.CreatedAt.String(),
		}
		newMessages = append(newMessages, newMessage)
		SC.ReceivedMessages[idUser] = append(SC.ReceivedMessages[idUser], message.ID)
	}
	log.Printf("Get %d new messages\n", len(newMessages))
	return &newMessages, nil
}
