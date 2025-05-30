package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type AdminSelectOldOrNewMsgsPage struct {
	*Page
}

func NewAdminSelectOldOrNewMsgsPage(injector *injector.AppInjector) *AdminSelectOldOrNewMsgsPage {
	p := &AdminSelectOldOrNewMsgsPage{
		Page: NewPage(injector),
	}
	p.setCurrenPage(p)

	return p
}

func (this *AdminSelectOldOrNewMsgsPage) Name() string {
	return consts.ADMIN_OLD_OR_NEW_MSG_PAGE
}

func (this *AdminSelectOldOrNewMsgsPage) AllowedOnlyCommands() bool {
	return true
}

func (this *AdminSelectOldOrNewMsgsPage) RespText(update tgbotapi.Update) string {
	return "Do you want to see only new messages?"
}

func (this *AdminSelectOldOrNewMsgsPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.AdminOldOrNewMessagesPageKeyboard
}

func (this *AdminSelectOldOrNewMsgsPage) ActionOnDestroy(update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	this.injector.AdminQueryBuilder.SetOldOrNewQueryMsg(
		this.UserName(update),
		update.CallbackData(),
	)
}

var _ models.IPageWithActionOnDestroy = (*AdminSelectOldOrNewMsgsPage)(nil)
var _ models.IPageWithKeyboard = (*AdminSelectOldOrNewMsgsPage)(nil)
