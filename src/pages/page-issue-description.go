package pages

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type IssueDescriptionPage struct {
	*Page
	errResp string
}

func NewIssueDescriptionPage(db models.IDatabase, adminQueryBuilder *query_builder.AdminQueryBuilder) *IssueDescriptionPage {
	return &IssueDescriptionPage{
		Page:    NewPage(db, adminQueryBuilder),
		errResp: "",
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

func (this *IssueDescriptionPage) Action(update tgbotapi.Update) {
	dbMsg := models.JsonClientMsg{
		UserName: this.UserName(update),
		Text:     update.Message.Text,
	}

	err := this.db.Tables().Messages.AddMessage(dbMsg)
	if err != nil {
		log.Println("[IssueDescriptionPage_Action] AddMessage err ==> ", err)
		this.errResp = "Error on server side trying to save your message. Try to contact support directly: https://t.me/eobuhow."
	} else {
		this.errResp = ""
	}
}

func (this *IssueDescriptionPage) NextPage(update tgbotapi.Update, isAdmin bool) models.IPage {
	if this.errResp != "" {
		return this
	}

	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case consts.COLLABORATE:
			return NewPartnershipPage(this.db, this.adminQueryBuilder)
		case consts.INTEGRATE:
			return NewIntegrationPage(this.db, this.adminQueryBuilder)
		case consts.SUPPORT:
			return NewSupportPage(this.db, this.adminQueryBuilder)
		case consts.OTHER:
			return NewOtherPage(this.db, this.adminQueryBuilder)
		case consts.DESCRIBE_ISSUE:
			return NewIssueDescriptionPage(this.db, this.adminQueryBuilder)

		case consts.SHOW_MESSAGES:
			return NewAdminSelectOldOrNewMsgsPage(this.db, this.adminQueryBuilder)
		case consts.CHECK_LINKS:
			return NewAdminLinksPage(this.db, this.adminQueryBuilder)
		case consts.SELECT_NUMBER_OF_MESSAGES:
			return NewAdminSelectMsgCountPage(this.db, this.adminQueryBuilder)
		case consts.SHOW_ALL_OR_NEW_MESSAGES:
			return NewAdminSelectOldOrNewMsgsPage(this.db, this.adminQueryBuilder)
		case consts.DELETE_MESSAGES:
			return NewAdminDeleteMsgCountPage(this.db, this.adminQueryBuilder)

		case consts.SHOW_ALL_MESSAGES, consts.SHOW_NEW_MESSAGES:
			return NewAdminSelectMsgCountPage(this.db, this.adminQueryBuilder)

		case consts.BACK_TO_START:
			if isAdmin {
				return NewAdminStartPage(this.db, this.adminQueryBuilder)
			} else {
				return NewStartPage(this.db, this.adminQueryBuilder)
			}
		default:
			if isAdmin {
				return NewAdminStartPage(this.db, this.adminQueryBuilder)
			} else {
				return NewStartPage(this.db, this.adminQueryBuilder)
			}
		}
	}

	return NewThanksPage(this.db, this.adminQueryBuilder)
}

var _ models.IPage = (*IssueDescriptionPage)(nil)
