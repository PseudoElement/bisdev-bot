package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type LinksForAdminPage struct {
	*Page
}

func NewAdminLinksPage(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *LinksForAdminPage {
	p := &LinksForAdminPage{
		Page: NewPage(db, bot, adminQueryBuilder),
	}
	p.setCurrenPage(p)

	return p
}

func (this *LinksForAdminPage) Name() string {
	return consts.ADMIN_LINKS_PAGE
}

func (this *LinksForAdminPage) RespText(update tgbotapi.Update) string {
	return `🔗 Rubic links:`
}

func (this *LinksForAdminPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.AdminLinksPageKeyboard
}

var _ models.IPageWithKeyboard = (*LinksForAdminPage)(nil)
