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

type UserOpenPage struct {
	UserName   string
	Initials   string
	OpenedPage string
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

func (this *Admin) NotSetChatID() bool {
	return this.ChatId == 0
}
