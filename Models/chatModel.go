package model

type Chat struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Users []*User `json:"users"`
}
