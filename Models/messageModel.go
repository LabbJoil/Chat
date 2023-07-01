package model

type Message struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Sender    *User  `json:"sender"`
	CreatedAt string `json:"createdAt"`
	IDChat    string `json:"idChat"`
	Error     string `json:"error"`
}
