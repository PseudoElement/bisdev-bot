package pages

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type StartPage struct {
	*AbstrUserPage
}

func NewStartPage(injector *injector.AppInjector) *StartPage {
	p := &StartPage{
		AbstrUserPage: NewAbstrUserPage(injector),
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
		`Hey %s! 👋 Thanks for reaching out to Rubic — your universal DeFi aggregator for finding the best rates.

To help you better, could you let me know what your request is about?

Please choose one of the options below:
`,
		this.UserName(update),
	)

	return resp
}

func (this *StartPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.StartPageKeyboard
}

var _ models.IUserPage = (*StartPage)(nil)
var _ models.IPageWithKeyboard = (*StartPage)(nil)
