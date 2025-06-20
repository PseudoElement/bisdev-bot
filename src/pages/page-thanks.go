package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type ThanksPage struct {
	*Page
}

func NewThanksPage(injector *injector.AppInjector) *ThanksPage {
	p := &ThanksPage{
		Page: NewPage(injector),
	}
	p.setCurrenPage(p)

	return p
}

func (this *ThanksPage) Name() string {
	return consts.THANKS_PAGE
}

func (this *ThanksPage) IsSelectablePage() bool {
	return false
}

func (this *ThanksPage) AllowedOnlyCommands() bool {
	return true
}

func (this *ThanksPage) RespText(update tgbotapi.Update) string {
	return "Thanks for reaching us! We'll process your request as soon as possible."
}

func (this *ThanksPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.BackToStartKeyBoard
}

var _ models.IPage = (*ThanksPage)(nil)
var _ models.IUserPage = (*ThanksPage)(nil)
var _ models.IPageWithKeyboard = (*ThanksPage)(nil)
