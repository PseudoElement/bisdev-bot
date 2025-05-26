package pages

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type Page struct {
	db                models.IDatabase
	adminQueryBuilder *query_builder.AdminQueryBuilder
	errResp           string
}

func NewPage(db models.IDatabase, adminQueryBuilder *query_builder.AdminQueryBuilder) *Page {
	return &Page{
		db:                db,
		adminQueryBuilder: adminQueryBuilder,
		errResp:           "",
	}
}

func (this *Page) AllowedOnlyCommands() bool {
	return false
}

func (this *Page) UserName(update tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.From.UserName
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.UserName
	}

	log.Println("[Page_UserName()] unexpected empty UserName ==> ", update)

	return ""
}

func (this *Page) NextPage(update tgbotapi.Update, isAdmin bool) models.IPage {
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
