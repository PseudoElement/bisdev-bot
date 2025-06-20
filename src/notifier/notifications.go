package notifier

type NotificationNewMessage struct {
	FromUserName string
	FromInitials string
	Text         string
	FileID       string
	FileType     string
}

type NotificationBlockUser struct {
	BlockedUserName string
	AdminUserName   string
}

type NotificationUnblockUser struct {
	UnblockedUserName string
	AdminUserName     string
}

type NotificationUserOpenPage struct {
	FromUserName string
	FromInitials string
	OpenedPage   string
}
