package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type AdminInfoAfterDeletionPage struct {
	*Page
}

func NewAdminInfoAfterDeletionPage(injector *injector.AppInjector) *AdminInfoAfterDeletionPage {
	p := &AdminInfoAfterDeletionPage{
		Page: NewPage(injector),
	}
	p.setCurrenPage(p)

	return p
}

func (this *AdminInfoAfterDeletionPage) Name() string {
	return consts.ADMIN_INFO_AFTER_MSG_DELETION_PAGE
}

func (this *AdminInfoAfterDeletionPage) AllowedOnlyCommands() bool {
	return true
}

func (this *AdminInfoAfterDeletionPage) RespText(update tgbotapi.Update) string {
	return "Messages were successfully deleted."
}

func (this *AdminInfoAfterDeletionPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.AdminInfoAfterDeletionMsgPageKeyboard
}

var _ models.IPageWithKeyboard = (*AdminInfoAfterDeletionPage)(nil)
