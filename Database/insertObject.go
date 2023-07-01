package databasePostgreSQL

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func (DBD *DatabaseDevelop) InserNewMessage(id string, senderID string, idConversation string, content string) error {
	queryString := fmt.Sprintf("INSERT INTO messages (id, content, sender, conversationid) VALUES ('%s', '%s', '%s', '%s');", id, content, senderID, idConversation)

	result := DBD.DB.Exec(queryString)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	log.Println("Message add")
	return nil
}
