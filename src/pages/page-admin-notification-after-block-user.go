package pages

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type NotificationAfterBlockUserPage struct {
	*Page
}

func NewNotificationAfterBlockUserPage(injector *injector.AppInjector) *NotificationAfterBlockUserPage {
	p := &NotificationAfterBlockUserPage{
		Page: NewPage(injector),
	}
	p.setCurrenPage(p)

	return p
}

func (this *NotificationAfterBlockUserPage) Name() string {
	return consts.ADMIN_NOTE_AFTER_USER_BLOCK_PAGE
}

func (this *NotificationAfterBlockUserPage) AllowedOnlyCommands() bool {
	return true
}

func (this *NotificationAfterBlockUserPage) RespText(update tgbotapi.Update) string {
	return fmt.Sprintf("User %s was blocked successfully.", this.TextFromClient(update))
}

func (this *NotificationAfterBlockUserPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.BackToStartKeyBoard
}

var _ models.IPageWithKeyboard = (*NotificationAfterBlockUserPage)(nil)
