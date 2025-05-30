package pages

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type NotificationAfterBlockUserPage struct {
	*Page
}

func NewNotificationAfterBlockUserPage(
	db models.IDatabase,
	bot *tgbotapi.BotAPI,
	adminQueryBuilder *query_builder.AdminQueryBuilder,
) *NotificationAfterBlockUserPage {
	p := &NotificationAfterBlockUserPage{
		Page: NewPage(db, bot, adminQueryBuilder),
	}
	p.setCurrenPage(p)

	return p
}

func (this *NotificationAfterBlockUserPage) Name() string {
	return consts.ADMIN_NOTE_AFTER_USER_BLOCK_PAGE
}

func (this *NotificationAfterBlockUserPage) AllowedOnlyCommands() bool {
	return true
}

func (this *NotificationAfterBlockUserPage) RespText(update tgbotapi.Update) string {
	return fmt.Sprintf("User %s was blocked successfully.", this.TextFromClient(update))
}

func (this *NotificationAfterBlockUserPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.BackToStartKeyBoard
}

var _ models.IPageWithKeyboard = (*NotificationAfterBlockUserPage)(nil)
