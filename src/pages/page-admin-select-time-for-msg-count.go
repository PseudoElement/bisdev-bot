package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type AdminSelectTimeForMsgCountPage struct {
	*Page
}

func NewAdminSelectTimeForMsgCountPage(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *AdminSelectTimeForMsgCountPage {
	p := &AdminSelectTimeForMsgCountPage{
		Page: NewPage(db, bot, adminQueryBuilder),
	}
	p.setCurrenPage(p)

	return p
}

func (this *AdminSelectTimeForMsgCountPage) Name() string {
	return consts.ADMIN_SELECT_TIME_FOR_MSG_COUNT_PAGE
}

func (this *AdminSelectTimeForMsgCountPage) AllowedOnlyCommands() bool {
	return true
}

func (this *AdminSelectTimeForMsgCountPage) RespText(update tgbotapi.Update) string {
	return "Select time interval users sent messages since:"
}

func (this *AdminSelectTimeForMsgCountPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.AdminSelectTimeForMsgCountPageKeyboard
}

var _ models.IPageWithKeyboard = (*AdminSelectTimeForMsgCountPage)(nil)
