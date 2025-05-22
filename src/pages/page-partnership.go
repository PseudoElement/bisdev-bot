package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type PartnershipPage struct {
	*Page
}

func NewPartnershipPage() *PartnershipPage {
	return &PartnershipPage{
		Page: NewPage(),
	}
}

func (this *PartnershipPage) Name() string {
	return consts.PARTNERSHIP_PAGE
}

func (this *PartnershipPage) RespText(update tgbotapi.Update) string {
	return `Awesome!🙌 Let's explore a potential collaboration.

Can you share the following:
- Project name:
- Website:
- Your role:
- Your main goal with us? (integration / liquidity aggregation / mutual routing / co-marketing / other)

Once you're done, I’ll share this with our BD team and we’ll follow up fast.  `
}

var _ models.IPage = (*PartnershipPage)(nil)
