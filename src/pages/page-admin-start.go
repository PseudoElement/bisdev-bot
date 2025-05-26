package pages

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type AdminStartPage struct {
	*Page
}

func NewAdminStartPage(db models.IDatabase, adminQueryBuilder *query_builder.AdminQueryBuilder) *AdminStartPage {
	return &AdminStartPage{
		Page: NewPage(db, adminQueryBuilder),
	}
}

func (this *AdminStartPage) Name() string {
	return consts.ADMIN_START_PAGE
}

func (this *AdminStartPage) AllowedOnlyCommands() bool {
	return true
}

func (this *AdminStartPage) RespText(update tgbotapi.Update) string {
	return fmt.Sprintf("Admin mode for **%s** activated. Select option:", this.UserName(update))
}

func (this *AdminStartPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.AdminStartPageKeyboard
}

var _ models.IPageWithKeyboard = (*AdminStartPage)(nil)
