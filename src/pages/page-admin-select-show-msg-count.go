package pages

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
)

type AdminSelectMsgCountPage struct {
	*Page
}

func NewAdminSelectMsgCountPage(injector *injector.AppInjector) *AdminSelectMsgCountPage {
	p := &AdminSelectMsgCountPage{
		Page: NewPage(injector),
	}
	p.setCurrenPage(p)

	return p
}

func (this *AdminSelectMsgCountPage) Name() string {
	return consts.ADMIN_SELECT_MSG_COUNT_PAGE
}

func (this *AdminSelectMsgCountPage) RespText(update tgbotapi.Update) string {
	if this.errResp != "" {
		return this.errResp
	}
	return "Input number of latest messages, you want to see (max number per request is 20):"
}

func (this *AdminSelectMsgCountPage) ActionOnDestroy(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	count, err := strconv.Atoi(this.TextFromClient(update))
	if err != nil {
		this.setErrorResp(fmt.Sprintf("%v is invalid number of messages.\n", this.TextFromClient(update)))
		return
	}
	if count <= 0 {
		this.setErrorResp("I think it's a joke. Try to use bigger number.")
		return
	}

	this.setErrorResp("")
	this.injector.AdminQueryBuilder.SetCountOfQueryMsg(
		this.UserName(update),
		count,
	)
}

func (this *AdminSelectMsgCountPage) NextPage(update tgbotapi.Update, isAdmin bool) models.IPage {
	if update.CallbackData() == consts.BACK_TO_START {
		return NewAdminStartPage(this.injector)
	} else if this.errResp != "" {
		return this
	} else {
		return NewAdminListOfMessagesPage(this.injector)
	}
}

func (this *AdminSelectMsgCountPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.BackToStartKeyBoard
}

var _ models.IPage = (*AdminSelectMsgCountPage)(nil)
var _ models.IPageWithKeyboard = (*AdminSelectMsgCountPage)(nil)
var _ models.IPageWithActionOnDestroy = (*AdminSelectMsgCountPage)(nil)
