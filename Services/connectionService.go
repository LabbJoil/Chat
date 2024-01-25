package service

import (
	database "github.com/LabbJoil/Chat/Database"
	databasemodels "github.com/LabbJoil/Chat/Models/DatabaseModels"
	"github.com/spf13/viper"
)

type ChatInteraction struct {
	DBConnection     database.DatabaseDevelop
	ReceivedMessages map[string][]string
}

func (SC *ChatInteraction) ConnectDB() error {
	DB := database.DatabaseDevelop{}
	DBConfig := databasemodels.DBConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		Password: viper.GetString("db.password"),
	}

	if err := DB.DBConnect(DBConfig); err != nil {
		return err
	}
	if err := DB.AutoMigrate(); err != nil {
		return err
	}
	SC.DBConnection = DB
	return nil
}

func IsUserInSlice(userIDToCheck string, users []*databasemodels.User) bool {
	for _, user := range users {
		if user.ID == userIDToCheck {
			return true
		}
	}
	return false
}
