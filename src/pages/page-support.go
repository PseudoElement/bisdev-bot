package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type SupportPage struct {
	*AbstrUserInputPage
}

func NewSupportPage(injector *injector.AppInjector) *SupportPage {
	p := &SupportPage{
		AbstrUserInputPage: NewAbstrUserInputPage(injector),
	}
	p.setCurrenPage(p)

	return p
}

func (this *SupportPage) Name() string {
	return consts.SUPPORT_PAGE
}

func (this *SupportPage) AllowedOnlyMessages() bool {
	return false
}

func (this *SupportPage) RespText(update tgbotapi.Update) string {
	if this.errResp != "" {
		return this.errResp
	}

	return `Sorry to hear you're having trouble 😔 Let's sort it out.
Please provide:
- A short description of the issue
- Tx Hash (if applicable)
- Network and tokens involved
- Device & browser (if on web)
- Screenshot (one image max per request)

🔧 For faster assistance, you can also contact our support team directly: https://t.me/RubicSupportBot.
Or just describe your issue here — I'll forward it to our tech team.`
}

func (this *SupportPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	if this.errResp != "" {
		return keyboards.BackToStartKeyBoard
	}
	return keyboards.SupportPageKeyboard
}

var _ models.IPage = (*SupportPage)(nil)
var _ models.IUserPage = (*SupportPage)(nil)
var _ models.IPageWithKeyboard = (*SupportPage)(nil)
var _ models.IPageWithActionOnDestroy = (*SupportPage)(nil)
