package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type SupportPage struct {
	*Page
}

func NewSupportPage() *SupportPage {
	return &SupportPage{
		Page: NewPage(),
	}
}

func (this *SupportPage) Name() string {
	return consts.SUPPORT_PAGE
}

func (this *SupportPage) RespText(update tgbotapi.Update) string {
	return `Sorry to hear you're having trouble 😔. Let me help.
Can you please provide:
- A short description of the issue
- Tx Hash (if relevant)
- Network / Tokens involved
- Device & browser (if on web)
- Screenshot (if possible)

🔧 For faster help, feel free to head to our support Telegram: https://t.me/eobuhow
Or describe your problem here — I’ll log this and escalate it to our tech support team.`
}

func (this *SupportPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.SupportPageKeyboard
}

var _ models.IPageWithKeyboard = (*SupportPage)(nil)
