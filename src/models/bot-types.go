package models

type JsonClientMsg struct {
	UserName string `json:"user_name"`
	Text     string `json:"text"`
}

type MessagesReq struct {
	Count   int
	NewOnly bool
}
