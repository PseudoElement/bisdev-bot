package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type AdminLinksPage struct {
	*Page
}

func NewAdminLinksPage(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *AdminLinksPage {
	p := &AdminLinksPage{
		Page: NewPage(db, bot, adminQueryBuilder),
	}
	p.setCurrenPage(p)

	return p
}

func (this *AdminLinksPage) AllowedOnlyCommands() bool {
	return true
}

func (this *AdminLinksPage) Name() string {
	return consts.ADMIN_LINKS_PAGE
}

func (this *AdminLinksPage) RespText(update tgbotapi.Update) string {
	return `🔗 Rubic links:`
}

func (this *AdminLinksPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.AdminLinksPageKeyboard
}

var _ models.IPageWithKeyboard = (*AdminLinksPage)(nil)
