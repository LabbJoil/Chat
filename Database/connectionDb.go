package databasePostgreSQL

import (
	"fmt"
	"log"

	databasemodels "github.com/LabbJoil/Chat/Models/DatabaseModels"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseDevelop struct {
	DB *gorm.DB
}

func (DBD *DatabaseDevelop) DBConnect(config databasemodels.DBConfig) error {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", config.Host, config.Username, config.Password, config.DBName, config.Port)

	dbAvailable, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return err
	}
	DBD.DB = dbAvailable
	return nil
}

func (DBD *DatabaseDevelop) AutoMigrate() error {

	if err := DBD.DB.AutoMigrate(&databasemodels.User{}); err != nil {
		log.Fatal(err)
		return err
	}
	if err := DBD.DB.AutoMigrate(&databasemodels.Conversation{}); err != nil {
		log.Fatal(err)
		return err
	}
	if err := DBD.DB.AutoMigrate(&databasemodels.ConversationUser{}); err != nil {
		log.Fatal(err)
		return err
	}
	if err := DBD.DB.AutoMigrate(&databasemodels.Message{}); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
