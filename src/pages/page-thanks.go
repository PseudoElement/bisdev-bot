package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type ThanksPage struct {
	*Page
}

func NewThanksPage(db models.IDatabase, adminQueryBuilder *query_builder.AdminQueryBuilder) *ThanksPage {
	return &ThanksPage{
		Page: NewPage(db, adminQueryBuilder),
	}
}

func (this *ThanksPage) Name() string {
	return consts.THANKS_PAGE
}

func (this *ThanksPage) AllowedOnlyCommands() bool {
	return true
}

func (this *ThanksPage) RespText(update tgbotapi.Update) string {
	return "Thanks for reaching us! We'll process your request as soon as possible."
}

func (this *ThanksPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.ThanksPageKeyboard
}

var _ models.IPageWithKeyboard = (*ThanksPage)(nil)
