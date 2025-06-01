package models

type UserMsgFromClient struct {
	UserId    int64
	UserName  string
	Initials  string
	Text      string
	FileType  string
	FileID    string
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
	FileID   string
	FileType string
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
