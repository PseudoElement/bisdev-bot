package models

type UserMsgFromClient struct {
	UserId    int64
	UserName  string
	Initials  string
	Text      string
	BlobType  string
	Blob      []byte
	CreatedAt string
}

type UserMsgWithFileID struct {
	*UserMsgFromClient
	FileID string
}

type MsgFromClientForLog struct {
	UserName string
	Initials string
	Text     string
	BlobLen  int
	BlobType string
}

type MessagesReq struct {
	Count   int
	NewOnly bool
}

type Admin struct {
	ChatId             int64
	UserName           string
	UserId             int64
	IsListenToNotifier bool
}
