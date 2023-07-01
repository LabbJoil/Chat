package databasePostgreSQL

import (
	"fmt"
	"log"

	model "github.com/LabbJoil/Chat/Models"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseDevelop struct {
	DB *gorm.DB
}

func (DBD *DatabaseDevelop) DBConnect(config model.DBConfig) error {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", config.Host, config.Username, config.Password, config.DBName, config.Port)

	dbAvailable, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return err
	}
	DBD.DB = dbAvailable
	return nil
}

func (DBD *DatabaseDevelop) TableExists() error {

	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id        VARCHAR(50) PRIMARY KEY,
		username  VARCHAR(50) not null,
		email     VARCHAR(100) not null
	)`

	createConversationsTable := `
	CREATE TABLE IF NOT EXISTS conversations (
		id        VARCHAR(50) PRIMARY KEY,
		name      VARCHAR(100),
		createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
		members   VARCHAR(1000)[]
	)`

	createMessageTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id                 VARCHAR(50) PRIMARY KEY,
		content            TEXT not null,
		sender	           VARCHAR(50) not null REFERENCES users (id) ON DELETE CASCADE,
		createdAt          TIMESTAMP NOT NULL DEFAULT NOW(),
		conversationId	   VARCHAR(50) not null REFERENCES conversations (id) ON DELETE CASCADE
	)`

	result := DBD.DB.Exec(createUserTable)
	if result.Error != nil {
		log.Fatal(result.Error)
		return result.Error
	}

	result = DBD.DB.Exec(createConversationsTable)
	if result.Error != nil {
		log.Fatal(result.Error)
		return result.Error
	}

	result = DBD.DB.Exec(createMessageTable)
	if result.Error != nil {
		log.Fatal(result.Error)
		return result.Error
	}

	return nil
}
