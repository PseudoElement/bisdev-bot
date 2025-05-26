package pages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type Page struct{}

func NewPage() *Page {
	return &Page{}
}

func (this *Page) AllowedOnlyCommands() bool {
	return false
}

func (this *Page) NextPage(update tgbotapi.Update, isAdmin bool) models.IPage {
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case consts.COLLABORATE:
			return NewPartnershipPage()
		case consts.INTEGRATE:
			return NewIntegrationPage()
		case consts.SUPPORT:
			return NewSupportPage()
		case consts.OTHER:
			return NewOtherPage()
		case consts.BACK_TO_START:
			if isAdmin {
				return NewAdminStartPage()
			} else {
				return NewStartPage()
			}
		case consts.DESCRIBE_ISSUE:
			return NewIssueDescriptionPage()
		default:
			return NewStartPage()
		}
	}

	return NewThanksPage()
}
