package pages

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type StartPage struct {
	*Page
}

func NewStartPage() *StartPage {
	return &StartPage{
		Page: NewPage(),
	}
}

func (this *StartPage) Name() string {
	return consts.START_PAGE
}

func (this *StartPage) AllowedOnlyCommands() bool {
	return true
}

func (this *StartPage) RespText(update tgbotapi.Update) string {
	var userName string
	if update.Message != nil {
		userName = update.Message.From.UserName
	}
	if update.CallbackQuery != nil {
		userName = update.CallbackQuery.From.UserName
	}

	resp := fmt.Sprintf(
		`Hey %s! 👋 Thanks for reaching out to Rubic — your universal DeFi aggregator for finding the best ratestraders. 
To help you better, can you tell me what your request is about?

Please choose one of the options below:
`,
		userName,
	)

	return resp
}

func (this *StartPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return startPageKeyboard
}

var _ models.IPageWithKeyboard = (*StartPage)(nil)
