package pages

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type AdminNotificationAfterUserUnblockPage struct {
	*Page
}

func NewAdminNotificationAfterUserUnblockPage(injector *injector.AppInjector) *AdminNotificationAfterUserUnblockPage {
	p := &AdminNotificationAfterUserUnblockPage{
		Page: NewPage(injector),
	}
	p.setCurrenPage(p)

	return p
}

func (this *AdminNotificationAfterUserUnblockPage) Name() string {
	return consts.ADMIN_NOTE_AFTER_USER_UNBLOCK_PAGE
}

func (this *AdminNotificationAfterUserUnblockPage) AllowedOnlyCommands() bool {
	return true
}

func (this *AdminNotificationAfterUserUnblockPage) RespText(update tgbotapi.Update) string {
	return fmt.Sprintf("User %s was unblocked successfully.", this.TextFromClient(update))
}

func (this *AdminNotificationAfterUserUnblockPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.BackToStartKeyBoard
}

var _ models.IPageWithKeyboard = (*AdminNotificationAfterUserUnblockPage)(nil)
