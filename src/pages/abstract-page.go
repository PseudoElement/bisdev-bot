package pages

import (
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type Page struct {
	injector *injector.AppInjector
	// db                models.IDatabase
	// adminQueryBuilder *query_builder.AdminQueryBuilder
	bot *tgbotapi.BotAPI
	// should be set in child structs
	currPage    models.IPage
	errResp     string
	warningResp string
}

func NewPage(injector *injector.AppInjector) *Page {
	return &Page{
		// db:                db,
		// adminQueryBuilder: adminQueryBuilder,
		injector:    injector,
		bot:         injector.Bot,
		currPage:    nil,
		errResp:     "",
		warningResp: "",
	}
}

func (this *Page) AllowedOnlyCommands() bool {
	return false
}

func (this *Page) AllowedOnlyMessages() bool {
	return false
}

func (this *Page) Bot() *tgbotapi.BotAPI {
	return this.bot
}

func (this *Page) CurrPage() models.IPage {
	return this.currPage
}

func (this *Page) UserName(update tgbotapi.Update) string {
	if update.Message != nil && update.Message.From.UserName != "" {
		return update.Message.From.UserName
	}
	if update.CallbackQuery != nil && update.CallbackQuery.From.UserName != "" {
		return update.CallbackQuery.From.UserName
	}

	log.Printf("[Page_UserName()] unexpected empty UserName ==> %+v", update.Message)

	return strconv.Itoa(int(update.Message.From.ID))
}

func (this *Page) UserID(update tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.From.ID
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}

	log.Println("[Page_UserID()] unexpected 0 UserID ==> ", update.Message)

	return 0
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

// @REDACTOR keep admins list in store
func (this *Page) IsUserAdmin(userName string) bool {
	userId := this.injector.Db.Tables().Messages.GetUserId(userName)

	adminsString, ok := os.LookupEnv("ADMINS")
	if !ok {
		return false
	}
	admins := strings.Split(adminsString, " ")

	for _, adminId := range admins {
		if strconv.Itoa(int(userId)) == adminId {
			return true
		}
	}

	return false
}

func (this *Page) NextPage(update tgbotapi.Update, isAdmin bool) models.IPage {
	// force return to start page if Comamnd BACK_TO_START selected
	if update.CallbackData() == consts.BACK_TO_START {
		if isAdmin {
			return NewAdminStartPage(this.injector)
		} else {
			return NewStartPage(this.injector)
		}
	}

	if this.errResp != "" || this.warningResp != "" {
		return this.CurrPage()
	}

	if update.CallbackQuery != nil {
		ttm := consts.TIME_TO_MIN
		switch update.CallbackData() {
		case consts.COLLABORATE:
			return NewPartnershipPage(this.injector)
		case consts.INTEGRATE:
			return NewIntegrationPage(this.injector)
		case consts.SUPPORT:
			return NewSupportPage(this.injector)
		case consts.OTHER:
			return NewOtherPage(this.injector)
		case consts.DESCRIBE_ISSUE:
			return NewIssueDescriptionPage(this.injector)

		case consts.SHOW_MESSAGES:
			return NewAdminSelectOldOrNewMsgsPage(this.injector)
		case consts.SHOW_MESSAGES_OF_SPECIFIC_USER, consts.DELETE_MESSAGES_OF_USER, consts.BLOCK_USER, consts.UNBLOCK_USER:
			return NewAdminInputUserNamePage(this.injector)
		case consts.CHECK_LINKS:
			return NewAdminLinksPage(this.injector)
		case consts.SHOW_ALL_OR_NEW_MESSAGES:
			return NewAdminSelectOldOrNewMsgsPage(this.injector)
		case consts.DELETE_MESSAGES:
			return NewAdminDeleteMsgCountPage(this.injector)
		case ttm.Mins_10, ttm.Mins_30, ttm.Hour_1, ttm.Hours_3, ttm.Hours_6, ttm.Hours_12, ttm.Day_1, ttm.Days_3, ttm.Week_1, ttm.Weeks_2, ttm.Month_1, ttm.Months_3:
			return NewAdminCountOfReceivedMsgsPage(this.injector)
		case consts.SHOW_ALL_MESSAGES, consts.SHOW_NEW_MESSAGES, consts.SELECT_NUMBER_OF_MESSAGES:
			return NewAdminSelectMsgCountPage(this.injector)
		case consts.SHOW_MESSAGES_COUNT_BY_TIME:
			return NewAdminSelectTimeForMsgCountPage(this.injector)

		default:
			if isAdmin {
				return NewAdminStartPage(this.injector)
			} else {
				return NewStartPage(this.injector)
			}
		}
	}

	return NewThanksPage(this.injector)
}

func (this *Page) setCurrenPage(page models.IPage) {
	this.currPage = page
}

func (this *Page) setErrorResp(err string) {
	this.errResp = err
}

func (this *Page) setWarningResp(warning string) {
	this.warningResp = warning
}
