package notifier

import "github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"

type Notifier struct {
	notificationChan chan any
}

func NewNotifier() *Notifier {
	n := &Notifier{
		notificationChan: make(chan any),
	}

	return n
}

func (this *Notifier) Chan() <-chan any {
	return this.notificationChan
}

func (this *Notifier) NotifyAdminsOnNewMsg(msg models.UserMsgFromClient) {
	this.notificationChan <- NotificationNewMessage{
		FromUserName: msg.UserName,
		FromInitials: msg.Initials,
		Text:         msg.Text,
		WithFiles:    len(msg.Blob) > 0,
	}
}

func (this *Notifier) NotifyAdminsOnBlockedUsers(userName string, adminName string) {
	this.notificationChan <- NotificationBlockUser{BlockedUserName: userName, AdminUserName: adminName}
}

func (this *Notifier) NotifyAdminsOnUnblockedUsers(userName string, adminName string) {
	this.notificationChan <- NotificationUnblockUser{UnblockedUserName: userName, AdminUserName: adminName}
}
