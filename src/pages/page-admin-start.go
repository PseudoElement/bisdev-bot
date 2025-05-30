package pages

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type AdminStartPage struct {
	*Page
}

func NewAdminStartPage(injector *injector.AppInjector) *AdminStartPage {
	p := &AdminStartPage{
		Page: NewPage(injector),
	}
	p.setCurrenPage(p)

	return p
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
