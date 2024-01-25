package remotemodels

type ChatMessage struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	IDSender  string `json:"sender"`
	CreatedAt string `json:"createdAt"`
	IDChat    string `json:"idChat"`
}
