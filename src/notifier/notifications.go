package notifier

type NotificationNewMessage struct {
	FromUserName string
	FromInitials string
	Text         string
	WithFiles    bool
}

type NotificationBlockUser struct {
	BlockedUserName string
	AdminUserName   string
}

type NotificationUnblockUser struct {
	UnblockedUserName string
	AdminUserName     string
}
