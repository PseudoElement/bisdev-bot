package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type AdminStartPage struct {
	*Page
}

func NewAdminStartPage() *AdminStartPage {
	return &AdminStartPage{
		Page: NewPage(),
	}
}

func (this *AdminStartPage) Name() string {
	return consts.ADMIN_START_PAGE
}

func (this *AdminStartPage) RespText(update tgbotapi.Update) string {
	return "Click button to see clients messages."
}

func (this *AdminStartPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.AdminStartPageKeyboard
}

var _ models.IPageWithKeyboard = (*AdminStartPage)(nil)
