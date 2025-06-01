package notifier

import (
	"log"

	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

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

func (this *Notifier) NotifyAdminsOnNewMsg(msg models.UserMsgWithFileID) {
	notification := NotificationNewMessage{
		FromUserName: msg.UserName,
		FromInitials: msg.Initials,
		Text:         msg.Text,
		FileID:       msg.FileID,
		FileType:     msg.BlobType,
	}
	log.Printf("[Notifier_NotifyAdminsOnNewMsg] newMsg_notification ==> %+v\n", notification)
	this.notificationChan <- notification
}

func (this *Notifier) NotifyAdminsOnBlockedUsers(userName string, adminName string) {
	notification := NotificationBlockUser{BlockedUserName: userName, AdminUserName: adminName}
	log.Printf("[Notifier_NotifyAdminsOnBlockedUsers] block_notification ==> %+v\n", notification)
	this.notificationChan <- notification
}

func (this *Notifier) NotifyAdminsOnUnblockedUsers(userName string, adminName string) {
	notification := NotificationUnblockUser{UnblockedUserName: userName, AdminUserName: adminName}
	log.Printf("[Notifier_NotifyAdminsOnUnblockedUsers] unblock_notification ==> %+v\n", notification)
	this.notificationChan <- notification
}
