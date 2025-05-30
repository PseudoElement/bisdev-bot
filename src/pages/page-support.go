package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type SupportPage struct {
	*AbstrUserInputPage
}

func NewSupportPage(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *SupportPage {
	basePage := NewPage(db, bot, adminQueryBuilder)
	p := &SupportPage{
		AbstrUserInputPage: NewAbstrUserInputPage(basePage),
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

	return `Sorry to hear you're having trouble 😔. Let me help.
Can you please provide:
- A short description of the issue
- Tx Hash (if relevant)
- Network / Tokens involved
- Device & browser (if on web)
- Screenshot (no more than 1 image per request)

🔧 For faster help, feel free to head to our support Telegram: https://t.me/RubicSupportBot.
Or describe your problem here — I’ll log this and escalate it to our tech support team.`
}

func (this *SupportPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	if this.errResp != "" {
		return keyboards.BackToStartKeyBoard
	}
	return keyboards.SupportPageKeyboard
}

var _ models.IPageWithKeyboard = (*SupportPage)(nil)
var _ models.IPageWithActionOnDestroy = (*SupportPage)(nil)
