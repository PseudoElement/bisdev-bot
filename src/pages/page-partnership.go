package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type PartnershipPage struct {
	*AbstrUserInputPage
}

func NewPartnershipPage(injector *injector.AppInjector) *PartnershipPage {
	p := &PartnershipPage{
		AbstrUserInputPage: NewAbstrUserInputPage(injector),
	}
	p.setCurrenPage(p)

	return p
}

func (this *PartnershipPage) Name() string {
	return consts.PARTNERSHIP_PAGE
}

func (this *PartnershipPage) RespText(update tgbotapi.Update) string {
	if this.errResp != "" {
		return this.errResp
	}

	return `Awesome!🙌 Let's explore a potential collaboration.

Could you please share the following:
- Project name
- Website
- Your role
- What are you looking for? (integration / liquidity aggregation / mutual routing / co-marketing / other)
- Screenshot (one image max per request)

Once submitted, I'll share this with our BD team — they'll get back to you shortly.`
}

func (this *PartnershipPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.BackToStartKeyBoard
}

var _ models.IPage = (*PartnershipPage)(nil)
var _ models.IUserPage = (*PartnershipPage)(nil)
var _ models.IPageWithKeyboard = (*PartnershipPage)(nil)
var _ models.IPageWithActionOnDestroy = (*PartnershipPage)(nil)
