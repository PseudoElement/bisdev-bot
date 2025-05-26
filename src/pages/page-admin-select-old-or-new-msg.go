package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type AdminSelectOldOrNewMsgsPage struct {
	*Page
}

func NewAdminSelectOldOrNewMsgsPage(db models.IDatabase, adminQueryBuilder *query_builder.AdminQueryBuilder) *AdminSelectOldOrNewMsgsPage {
	return &AdminSelectOldOrNewMsgsPage{
		Page: NewPage(db, adminQueryBuilder),
	}
}

func (this *AdminSelectOldOrNewMsgsPage) Name() string {
	return consts.ADMIN_OLD_OR_NEW_MSG_PAGE
}

func (this *AdminSelectOldOrNewMsgsPage) AllowedOnlyCommands() bool {
	return true
}

func (this *AdminSelectOldOrNewMsgsPage) RespText(update tgbotapi.Update) string {
	return "Do you want to see only new messages?"
}

func (this *AdminSelectOldOrNewMsgsPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.AdminOldOrNewMessagesPageKeyboard
}

func (this *AdminSelectOldOrNewMsgsPage) Action(update tgbotapi.Update) {
	this.adminQueryBuilder.SetOldOrNewQueryMsg(
		this.UserName(update),
		update.CallbackData(),
	)
}

var _ models.IPageWithAction = (*AdminSelectOldOrNewMsgsPage)(nil)
var _ models.IPageWithKeyboard = (*AdminSelectOldOrNewMsgsPage)(nil)
