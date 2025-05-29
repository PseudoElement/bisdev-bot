package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type IssueDescriptionPage struct {
	*AbstrUserInputPage
}

func NewIssueDescriptionPage(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *IssueDescriptionPage {
	basePage := NewPage(db, bot, adminQueryBuilder)
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

var _ models.IPageWithActionOnDestroy = (*IssueDescriptionPage)(nil)
