package models

type ClientMsg struct {
	UserName string
	Text     string
	Id       int64
	IsNew    bool
}

type JsonClientMsg struct {
	UserName string `json:"user_name"`
	Text     string `json:"text"`
}

type MessagesReq struct {
	Count   int
	NewOnly bool
}
