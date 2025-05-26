package pages

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type AdminDeleteMsgCountPage struct {
	*Page
	errRespMsg string
}

func NewAdminDeleteMsgCountPage(db models.IDatabase, adminQueryBuilder *query_builder.AdminQueryBuilder) *AdminDeleteMsgCountPage {
	return &AdminDeleteMsgCountPage{
		Page:       NewPage(db, adminQueryBuilder),
		errRespMsg: "",
	}
}

func (this *AdminDeleteMsgCountPage) Name() string {
	return consts.ADMIN_MSG_COUNT_PAGE
}

func (this *AdminDeleteMsgCountPage) RespText(update tgbotapi.Update) string {
	if this.errRespMsg != "" {
		return this.errRespMsg
	}
	return "Input number of latest messages, you want to delete (ex. 1, 10):"
}

// add "ALL" messages
func (this *AdminDeleteMsgCountPage) Action(update tgbotapi.Update) {
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
	go func() {
		err := this.db.Tables().Messages.DeleteMessages(count)
		if err != nil {
			log.Println("[AdminDeleteMsgCountPage_Action] err in DeleteMessages ==> ", err)
		}
	}()
}

func (this *AdminDeleteMsgCountPage) NextPage(update tgbotapi.Update, isAdmin bool) models.IPage {
	if this.errRespMsg != "" {
		return this
	} else {
		return NewAdminInfoAfterDeletionPage(this.db, this.adminQueryBuilder)
	}
}

var _ models.IPageWithAction = (*AdminDeleteMsgCountPage)(nil)
