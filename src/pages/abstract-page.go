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
	bot               *tgbotapi.BotAPI
	// should be set in child structs
	currPage models.IPage
	errResp  string
}

func NewPage(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *Page {
	return &Page{
		db:                db,
		adminQueryBuilder: adminQueryBuilder,
		bot:               bot,
		currPage:          nil,
		errResp:           "",
	}
}

func (this *Page) AllowedOnlyCommands() bool {
	return false
}

func (this *Page) Bot() *tgbotapi.BotAPI {
	return this.bot
}

func (this *Page) CurrPage() models.IPage {
	return this.currPage
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

func (this *Page) TextFromClient(update tgbotapi.Update) string {
	if update.Message != nil {
		if update.Message.Text != "" {
			return update.Message.Text
		}
		if update.Message.Caption != "" {
			return update.Message.Caption
		}
	}
	if update.CallbackQuery != nil {
		return update.CallbackData()
	}
	return ""
}

func (this *Page) NextPage(update tgbotapi.Update, isAdmin bool) models.IPage {
	if this.errResp != "" {
		return this.CurrPage()
	}

	if update.CallbackQuery != nil {
		ttm := consts.TIME_TO_MIN
		switch update.CallbackData() {
		case consts.COLLABORATE:
			return NewPartnershipPage(this.db, this.bot, this.adminQueryBuilder)
		case consts.INTEGRATE:
			return NewIntegrationPage(this.db, this.bot, this.adminQueryBuilder)
		case consts.SUPPORT:
			return NewSupportPage(this.db, this.bot, this.adminQueryBuilder)
		case consts.OTHER:
			return NewOtherPage(this.db, this.bot, this.adminQueryBuilder)
		case consts.DESCRIBE_ISSUE:
			return NewIssueDescriptionPage(this.db, this.bot, this.adminQueryBuilder)

		case consts.SHOW_MESSAGES:
			return NewAdminSelectOldOrNewMsgsPage(this.db, this.bot, this.adminQueryBuilder)
		case consts.SHOW_MESSAGES_OF_SPECIFIC_USER, consts.DELETE_MESSAGES_OF_USER:
			return NewAdminInputUserNamePage(this.db, this.bot, this.adminQueryBuilder)
		case consts.CHECK_LINKS:
			return NewAdminLinksPage(this.db, this.bot, this.adminQueryBuilder)
		case consts.SELECT_NUMBER_OF_MESSAGES:
			return NewAdminSelectMsgCountPage(this.db, this.bot, this.adminQueryBuilder)
		case consts.SHOW_ALL_OR_NEW_MESSAGES:
			return NewAdminSelectOldOrNewMsgsPage(this.db, this.bot, this.adminQueryBuilder)
		case consts.DELETE_MESSAGES:
			return NewAdminDeleteMsgCountPage(this.db, this.bot, this.adminQueryBuilder)
		case ttm.Mins_10, ttm.Mins_30, ttm.Hour_1, ttm.Hours_3, ttm.Hours_6, ttm.Hours_12, ttm.Day_1, ttm.Days_3, ttm.Week_1, ttm.Weeks_2, ttm.Month_1, ttm.Months_3:
			return NewAdminCountOfReceivedMsgsPage(this.db, this.bot, this.adminQueryBuilder)
		case consts.SHOW_ALL_MESSAGES, consts.SHOW_NEW_MESSAGES:
			return NewAdminSelectMsgCountPage(this.db, this.bot, this.adminQueryBuilder)
		case consts.SHOW_MESSAGES_COUNT_BY_TIME:
			return NewAdminSelectTimeForMsgCountPage(this.db, this.bot, this.adminQueryBuilder)

		case consts.BACK_TO_START:
			if isAdmin {
				return NewAdminStartPage(this.db, this.bot, this.adminQueryBuilder)
			} else {
				return NewStartPage(this.db, this.bot, this.adminQueryBuilder)
			}
		default:
			if isAdmin {
				return NewAdminStartPage(this.db, this.bot, this.adminQueryBuilder)
			} else {
				return NewStartPage(this.db, this.bot, this.adminQueryBuilder)
			}
		}
	}

	return NewThanksPage(this.db, this.bot, this.adminQueryBuilder)
}

func (this *Page) setCurrenPage(page models.IPage) {
	this.currPage = page
}

func (this *Page) setErrorResp(err string) {
	this.errResp = err
}
