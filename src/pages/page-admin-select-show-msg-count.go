package pages

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type AdminSelectMsgCountPage struct {
	*Page
	errRespMsg string
}

func NewAdminSelectMsgCountPage(db models.IDatabase, adminQueryBuilder *query_builder.AdminQueryBuilder) *AdminSelectMsgCountPage {
	return &AdminSelectMsgCountPage{
		Page:       NewPage(db, adminQueryBuilder),
		errRespMsg: "",
	}
}

func (this *AdminSelectMsgCountPage) Name() string {
	return consts.ADMIN_MSG_COUNT_PAGE
}

func (this *AdminSelectMsgCountPage) RespText(update tgbotapi.Update) string {
	if this.errRespMsg != "" {
		return this.errRespMsg
	}
	return "Input number of latest messages, you want to see (ex. 1, 10, all):"
}

// add "ALL" messages
func (this *AdminSelectMsgCountPage) Action(update tgbotapi.Update) {
	count, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		this.errRespMsg = fmt.Sprintf("%v is invalid number of messages.\n", update.Message.Text)
		return
	}
	if count == 0 {
		this.errRespMsg = "I think it's a joke. Try to use bigger number."
		return
	}

	this.errRespMsg = ""
	this.adminQueryBuilder.SetCountOfQueryMsg(
		this.UserName(update),
		count,
	)
}

func (this *AdminSelectMsgCountPage) NextPage(update tgbotapi.Update, isAdmin bool) models.IPage {
	if this.errRespMsg != "" {
		return this
	} else {
		return NewAdminListOfMessagesPage(this.db, this.adminQueryBuilder)
	}
}

var _ models.IPageWithAction = (*AdminSelectMsgCountPage)(nil)
