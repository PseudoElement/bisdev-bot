package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type OtherPage struct {
	*Page
}

func NewOtherPage() *OtherPage {
	return &OtherPage{
		Page: NewPage(),
	}
}

func (this *OtherPage) Name() string {
	return consts.OTHER_PAGE
}

func (this *OtherPage) RespText(update tgbotapi.Update) string {
	return `No problem! Please describe your request in a few words — I’ll make sure it reaches the right person on our team.
We aim to reply within 24 hours.`
}

var _ models.IPage = (*OtherPage)(nil)
