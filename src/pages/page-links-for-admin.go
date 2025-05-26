package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type LinksForAdminPage struct {
	*Page
}

func NewLinksForAdminPage() *LinksForAdminPage {
	return &LinksForAdminPage{
		Page: NewPage(),
	}
}

func (this *LinksForAdminPage) Name() string {
	return consts.LINKS_FOR_ADMIN_PAGE
}

func (this *LinksForAdminPage) RespText(update tgbotapi.Update) string {
	return `🔗 Rubic links:`
}

func (this *LinksForAdminPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.LinksForAdminPageKeyboard
}

var _ models.IPageWithKeyboard = (*LinksForAdminPage)(nil)
