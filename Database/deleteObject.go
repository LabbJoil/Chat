package databasePostgreSQL

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func (DBD *DatabaseDevelop) DeleteMessageById(idMessage string) error {
	queryString := fmt.Sprintf("Delete from messages where id = '%v';", idMessage)

	result := DBD.DB.Exec(queryString)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		newErr := fmt.Errorf("Any error. Message with id: %s not delete", idMessage)
		log.Println(newErr)
		return newErr
	}

	log.Printf("Message with id: %s delete \n", idMessage)
	return nil
}

func (DBD *DatabaseDevelop) DeleteChatByUserId(idUser string) error {
	queryString := fmt.Sprintf("Delete from conversations WHERE '%v' = ANY (members);", idUser)

	result := DBD.DB.Exec(queryString)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		newErr := fmt.Errorf("Any error. Chat with id: %s not delete", idUser)
		log.Println(newErr)
		return newErr
	}

	log.Printf("Chat with id: %s delete \n", idUser)
	return nil
}
