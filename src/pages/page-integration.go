package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type IntegrationPage struct {
	*Page
}

func NewIntegrationPage(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *IntegrationPage {
	p := &IntegrationPage{
		Page: NewPage(db, bot, adminQueryBuilder),
	}
	p.setCurrenPage(p)

	return p
}

func (this *IntegrationPage) Name() string {
	return consts.INTEGRATION_PAGE
}

func (this *IntegrationPage) AllowedOnlyCommands() bool {
	return true
}

func (this *IntegrationPage) RespText(update tgbotapi.Update) string {
	return `Here’s everything you need to get started:`
}

func (this *IntegrationPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.IntegrationPageKeyboard
}

var _ models.IPageWithKeyboard = (*IntegrationPage)(nil)
