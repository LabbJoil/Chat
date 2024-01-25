package databasePostgreSQL

import (
	databasemodels "github.com/LabbJoil/Chat/Models/DatabaseModels"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func (DBD *DatabaseDevelop) GetUser(idUser string) (*databasemodels.User, error) {
	var user databasemodels.User
	result := DBD.DB.Where("Id = ?", idUser).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (DBD *DatabaseDevelop) InsertNewUser(name, email string) (*databasemodels.User, error) {
	user := &databasemodels.User{
		ID:       "user_" + uuid.New().String(),
		Username: name,
		Email:    email,
	}
	result := DBD.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (DBD *DatabaseDevelop) DeleteUser(idUser string) error {
	user := &databasemodels.User{
		ID: idUser,
	}
	result := DBD.DB.Delete(user)
	return result.Error
}
