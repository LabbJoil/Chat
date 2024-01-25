package databasemodels

type User struct {
	ID       string `gorm:"primarykey"`
	Username string `gorm:"not null"`
	Email    string `gorm:"not null"`
}
