package models

type JsonMsgFromClient struct {
	UserName  string
	Initials  string
	Text      string
	ImageBlob []byte
}

type MsgFromClientForLog struct {
	UserName     string
	Initials     string
	Text         string
	ImageBlobLen int
}

type MessagesReq struct {
	Count   int
	NewOnly bool
}
