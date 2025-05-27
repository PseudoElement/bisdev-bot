package models

type JsonMsgFromClient struct {
	UserName  string `json:"user_name"`
	Text      string `json:"text"`
	ImageBlob []byte `json:"image_blob"`
}

type MessagesReq struct {
	Count   int
	NewOnly bool
}
