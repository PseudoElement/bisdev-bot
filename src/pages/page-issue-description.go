package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type IssueDescriptionPage struct {
	*Page
}

func NewIssueDescriptionPage() *IssueDescriptionPage {
	return &IssueDescriptionPage{
		Page: NewPage(),
	}
}

func (this *IssueDescriptionPage) Name() string {
	return consts.DESCRIBE_ISSUE
}

func (this *IssueDescriptionPage) RespText(update tgbotapi.Update) string {
	return `Can you please provide:
- A short description of the issue
- Tx Hash (if relevant)
- Network / Tokens involved
- Device & browser (if on web)
- Screenshot (if possible)`
}

var _ models.IPage = (*IssueDescriptionPage)(nil)
