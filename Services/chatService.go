package service

import (
	"fmt"
	"log"
	"time"

	database "github.com/LabbJoil/Chat/Database"
	model "github.com/LabbJoil/Chat/Models"
	"github.com/spf13/viper"
)

type ChatInteraction struct {
	DBConnection  database.DatabaseDevelop
	GetIdMessages []string
}

func (SC *ChatInteraction) ConnectDB() error {
	DB := database.DatabaseDevelop{}
	DBConfig := model.DBConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		Password: viper.GetString("db.password"),
	}

	if err := DB.DBConnect(DBConfig); err != nil {
		return err
	}
	if err := DB.TableExists(); err != nil {
		return err
	}
	SC.DBConnection = DB
	return nil
}

func (SC *ChatInteraction) GenerateMessage(senderID string, idConversation string, content string) (*model.Message, error) {
	membersChat, err := SC.DBConnection.GetConversationMembers(idConversation)
	if err != nil {
		return nil, err
	}

	if !FindElement(membersChat, senderID) {
		newErr := fmt.Errorf("Any error. User:%s not in chat", senderID)
		log.Println(newErr)
		return nil, newErr
	}
	senderUser, err := SC.DBConnection.GetUser(senderID)
	if err != nil {
		return nil, err
	}

	idMessage := GenerateMessageID(senderID, idConversation)
	err = SC.DBConnection.InserNewMessage(idMessage, senderID, idConversation, content)
	if err != nil {
		return nil, err
	}

	return &model.Message{
		ID:      idMessage,
		Content: content,
		Sender:  senderUser,
		IDChat:  idConversation,
		Error:   "nil",
	}, nil
}

func (SC *ChatInteraction) GetMessagesByUserId(idUser string) (*[]model.Message, error) {
	channelMessages := SC.DBConnection.GetMessages(append([]string{idUser}, SC.GetIdMessages...), "notRead")
	newMessages := []model.Message{}
	for message := range channelMessages {
		if message.Error != "nil" {
			return nil, fmt.Errorf(message.Error)
		}
		newMessages = append(newMessages, message)
		SC.GetIdMessages = append(SC.GetIdMessages, message.ID)
	}
	return &newMessages, nil
}

func (SC *ChatInteraction) GetAllMessagesByUserId(idUserFrom, idUserQuery, idMesssage, count string) ([]*model.Message, error) {
	channelMessages := SC.DBConnection.GetMessages([]string{idUserFrom, idUserQuery, idMesssage, count}, "allMessages")
	newMessages := []*model.Message{}
	for message := range channelMessages {
		if message.Error != "nil" {
			return nil, fmt.Errorf(message.Error)
		}
		newMessages = append(newMessages, &model.Message{
			ID:        message.ID,
			IDChat:    message.IDChat,
			Content:   message.Content,
			Sender:    message.Sender,
			CreatedAt: message.CreatedAt,
		})
	}
	return newMessages, nil
}

func (SC *ChatInteraction) GetUserChat(idUser string) ([]*model.Chat, error) {
	allUserChat, err := SC.DBConnection.GetAllChats(idUser)
	if err != nil {
		return nil, err
	}
	return allUserChat, nil
}

func (SC *ChatInteraction) DeleteMessageByMessage(idMessage string) error {
	if err := SC.DBConnection.DeleteMessageById(idMessage); err != nil {
		return err
	}
	return nil
}

func (SC *ChatInteraction) DeleteChatByUser(idUser string) error {
	if err := SC.DBConnection.DeleteChatByUserId(idUser); err != nil {
		return err
	}
	return nil
}

func GenerateMessageID(senderID, chatID string) string {
	currentTime := time.Now()
	messageID := fmt.Sprintf("%s_%s_%d", senderID, chatID, currentTime.UnixNano())
	return messageID
}

func FindElement(sliceString []string, element string) bool {
	for _, sliceElem := range sliceString {
		if element == sliceElem {
			return true
		}
	}
	return false
}
