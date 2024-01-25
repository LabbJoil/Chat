package remotemodels

type Chat struct {
	ID           string `json:"id"`
	Title        string `json:"name"`
	CountMembers int    `json:"countMembers"`
}
