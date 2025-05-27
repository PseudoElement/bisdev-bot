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
}

func NewAdminDeleteMsgCountPage(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *AdminDeleteMsgCountPage {
	p := &AdminDeleteMsgCountPage{
		Page: NewPage(db, bot, adminQueryBuilder),
	}
	p.setCurrenPage(p)

	return p
}

func (this *AdminDeleteMsgCountPage) Name() string {
	return consts.ADMIN_MSG_COUNT_PAGE
}

func (this *AdminDeleteMsgCountPage) RespText(update tgbotapi.Update) string {
	if this.errResp != "" {
		return this.errResp
	}
	return "Input number of latest messages, you want to delete (ex. 1, 10):"
}

// add "ALL" messages
func (this *AdminDeleteMsgCountPage) ActionOnDestroy(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	count, err := strconv.Atoi(this.TextFromClient(update))
	if err != nil {
		this.setErrorResp(fmt.Sprintf("%v is invalid number of messages.\n", this.TextFromClient(update)))
		return
	}
	if count == 0 {
		this.setErrorResp("I think it's a joke. Try to use bigger number.")
		return
	}

	this.setErrorResp("")
	go func() {
		err := this.db.Tables().Messages.DeleteMessages(count)
		if err != nil {
			log.Println("[AdminDeleteMsgCountPage_Action] err in DeleteMessages ==> ", err)
		}
	}()
}

func (this *AdminDeleteMsgCountPage) NextPage(update tgbotapi.Update, isAdmin bool) models.IPage {
	if this.errResp != "" {
		return this
	} else {
		return NewAdminInfoAfterDeletionPage(this.db, this.bot, this.adminQueryBuilder)
	}
}

var _ models.IPageWithActionOnDestroy = (*AdminDeleteMsgCountPage)(nil)
