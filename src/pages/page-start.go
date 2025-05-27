package pages

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type StartPage struct {
	*Page
}

func NewStartPage(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *StartPage {
	p := &StartPage{
		Page: NewPage(db, bot, adminQueryBuilder),
	}
	p.setCurrenPage(p)

	return p
}

func (this *StartPage) Name() string {
	return consts.START_PAGE
}

func (this *StartPage) AllowedOnlyCommands() bool {
	return true
}

func (this *StartPage) RespText(update tgbotapi.Update) string {
	resp := fmt.Sprintf(
		`Hey %s! 👋 Thanks for reaching out to Rubic — your universal DeFi aggregator for finding the best ratestraders. 
To help you better, can you tell me what your request is about?

Please choose one of the options below:
`,
		this.UserName(update),
	)

	return resp
}

func (this *StartPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.StartPageKeyboard
}

var _ models.IPageWithKeyboard = (*StartPage)(nil)
