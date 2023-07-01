package databasePostgreSQL

import (
	"fmt"
	"log"
	"strings"

	model "github.com/LabbJoil/Chat/Models"
	_ "github.com/lib/pq"
)

func (DBD *DatabaseDevelop) GetUser(idUser string) (*model.User, error) {
	queryString := fmt.Sprintf("Select * from users where id = '%v';", idUser)
	var user model.User

	result := DBD.DB.Raw(queryString).Scan(&user)
	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}
	log.Printf("Get user with id: %v\n", user.ID)
	return &user, nil
}

func (DBD *DatabaseDevelop) GetMessages(masFilterObject []string, typeQuery string) chan model.Message {
	messageChannel := make(chan model.Message)

	go func() {
		var queryString string
		switch typeQuery {
		case "notRead":
			queryString = fmt.Sprintf("Select messages.id, content, sender, conversations.id, messages.createdat from messages inner join conversations on messages.conversationid = conversations.id WHERE '%v' = ANY (members) AND messages.id not in ('%v');", masFilterObject[0], strings.Join(masFilterObject[1:], "', '"))

		case "allMessages":
			if masFilterObject[2] != "-" {
				queryString = fmt.Sprintf("Select messages.id, content, sender, conversations.id, messages.createdat from messages inner join conversations on messages.conversationid = conversations.id WHERE '%s'= messages.sender AND '%s' = ANY (members) AND '%s' = ANY(members) AND messages.createdat > (SELECT createdat FROM messages WHERE id = '%s') Order by messages.createdat Limit %s;", masFilterObject[0], masFilterObject[0], masFilterObject[1], masFilterObject[2], masFilterObject[3])
			} else {
				queryString = fmt.Sprintf("Select messages.id, content, sender, conversations.id, messages.createdat from messages inner join conversations on messages.conversationid = conversations.id WHERE '%s'= messages.sender AND '%s' = ANY (members) AND '%s' = ANY(members) Order by messages.createdat Limit '%s';", masFilterObject[0], masFilterObject[0], masFilterObject[1], masFilterObject[3])
			}

		}

		resultRows, err := DBD.DB.Raw(queryString).Rows()
		if err != nil {
			messageChannel <- model.Message{Error: err.Error()}
			log.Println(err.Error())
			close(messageChannel)
			return
		}
		defer resultRows.Close()
		var idMessage, content, senderId, idChat, createdat string

		for resultRows.Next() {
			err = resultRows.Scan(&idMessage, &content, &senderId, &idChat, &createdat)
			if err != nil {
				messageChannel <- model.Message{Error: err.Error()}
				log.Println(err.Error())
				close(messageChannel)
				return
			}

			senderUser, err := DBD.GetUser(senderId)
			if err != nil {
				messageChannel <- model.Message{Error: err.Error()}
				log.Println(err.Error())
				close(messageChannel)
				return
			}

			messageChannel <- model.Message{
				ID:        idMessage,
				Content:   content,
				Sender:    senderUser,
				IDChat:    idChat,
				CreatedAt: createdat,
				Error:     "nil",
			}
			log.Printf("Get message with id: %s\n", idMessage)
		}
		close(messageChannel)
	}()
	return messageChannel
}

func (DBD *DatabaseDevelop) GetAllChats(idUser string) ([]*model.Chat, error) {
	queryString := fmt.Sprintf("Select id, name, members from conversations where '%v' = ANY (members);", idUser)

	resultRows, err := DBD.DB.Raw(queryString).Rows()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resultRows.Close()

	var allChats []*model.Chat
	var idChat, nameChat, membersString string

	for resultRows.Next() {
		err = resultRows.Scan(&idChat, &nameChat, &membersString)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		membersChat := strings.Split(strings.Trim(string(membersString), "{}"), ",")
		var allChatUsers []*model.User

		for _, userId := range membersChat {
			user, err := DBD.GetUser(userId)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			allChatUsers = append(allChatUsers, user)
		}
		chat := model.Chat{
			ID:    idChat,
			Name:  nameChat,
			Users: allChatUsers,
		}
		allChats = append(allChats, &chat)
		log.Printf("Get chat with id: %v\n", chat.ID)
	}
	return allChats, nil
}

func (DBD *DatabaseDevelop) GetConversationMembers(idMessage string) ([]string, error) {
	queryString := fmt.Sprintf("Select members from conversations where id = '%v';", idMessage)
	var membersChatString string

	result := DBD.DB.Raw(queryString).Scan(&membersChatString)
	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}
	membersChat := strings.Split(strings.Trim(string(membersChatString), "{}"), ",")
	log.Printf("Get members with ids: %v\n", membersChat)
	return membersChat, nil
}
