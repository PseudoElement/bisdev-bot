package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type IntegrationPage struct {
	*Page
}

func NewIntegrationPage(injector *injector.AppInjector) *IntegrationPage {
	p := &IntegrationPage{
		Page: NewPage(injector),
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
	return `Here's everything you need to get started:`
}

func (this *IntegrationPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.IntegrationPageKeyboard
}

var _ models.IPageWithKeyboard = (*IntegrationPage)(nil)
