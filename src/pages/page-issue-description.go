package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type IssueDescriptionPage struct {
	*AbstrUserInputPage
}

func NewIssueDescriptionPage(injector *injector.AppInjector) *IssueDescriptionPage {
	basePage := NewPage(injector)
	p := &IssueDescriptionPage{
		AbstrUserInputPage: NewAbstrUserInputPage(basePage),
	}
	p.setCurrenPage(p)

	return p
}

func (this *IssueDescriptionPage) Name() string {
	return consts.DESCRIBE_ISSUE
}

func (this *IssueDescriptionPage) RespText(update tgbotapi.Update) string {
	if this.errResp != "" {
		return this.errResp
	}

	return `Can you please provide:
- A short description of the issue
- Tx Hash (if relevant)
- Network / Tokens involved
- Device & browser (if on web)
- Screenshot (no more than 1 image per request)`
}

func (this *IssueDescriptionPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.BackToStartKeyBoard
}

var _ models.IPage = (*IssueDescriptionPage)(nil)
var _ models.IPageWithKeyboard = (*IssueDescriptionPage)(nil)
var _ models.IPageWithActionOnDestroy = (*IssueDescriptionPage)(nil)
