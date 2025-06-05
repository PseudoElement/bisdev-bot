package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type OtherPage struct {
	*AbstrUserInputPage
}

func NewOtherPage(injector *injector.AppInjector) *OtherPage {
	basePage := NewPage(injector)
	p := &OtherPage{
		AbstrUserInputPage: NewAbstrUserInputPage(basePage),
	}
	p.setCurrenPage(p)

	return p
}

func (this *OtherPage) Name() string {
	return consts.OTHER_PAGE
}

func (this *OtherPage) RespText(update tgbotapi.Update) string {
	if this.errResp != "" {
		return this.errResp
	}

	return `No problem!
Just describe your request briefly — I'll make sure it reaches the right person on our team.
We aim to reply within 24 hours. (Please include no more than 1 image per request.)`
}

func (this *OtherPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.BackToStartKeyBoard
}

var _ models.IPage = (*OtherPage)(nil)
var _ models.IPageWithKeyboard = (*OtherPage)(nil)
var _ models.IPageWithActionOnDestroy = (*OtherPage)(nil)
