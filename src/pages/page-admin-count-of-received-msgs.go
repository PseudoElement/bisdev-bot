package pages

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/utils"
)

type AdminCountOfReceivedMsgsPage struct {
	*Page
	msgCount int
}

func NewAdminCountOfReceivedMsgsPage(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *AdminCountOfReceivedMsgsPage {
	p := &AdminCountOfReceivedMsgsPage{
		Page:     NewPage(db, bot, adminQueryBuilder),
		msgCount: 0,
	}
	p.setCurrenPage(p)

	return p
}

func (this *AdminCountOfReceivedMsgsPage) Name() string {
	return consts.ADMIN_RECEIVED_MSG_COUNT_PAGE
}

func (this *AdminCountOfReceivedMsgsPage) AllowedOnlyCommands() bool {
	return true
}

func (this *AdminCountOfReceivedMsgsPage) RespText(update tgbotapi.Update) string {
	if this.errResp != "" {
		return this.errResp
	}

	inputNum := this.TextFromClient(update)
	minsCountAgo, _ := strconv.Atoi(inputNum)
	timestamp := utils.GetSqlTimestampByMinutes(minsCountAgo, true)

	str := "messages"
	if this.msgCount == 1 {
		str = "message"
	}

	return fmt.Sprintf("Users sent %d %s since %s.", this.msgCount, str, timestamp)
}

func (this *AdminCountOfReceivedMsgsPage) ActionOnInit(update tgbotapi.Update) {
	inputNum := this.TextFromClient(update)
	minsCountAgo, err := strconv.Atoi(inputNum)
	if err != nil {
		this.setErrorResp(inputNum + "is invalid time.")
		return
	}

	sqlTimestamp := utils.GetSqlTimestampByMinutesUTC(minsCountAgo, true)
	msgCount, err := this.db.Tables().Messages.CheckMessagesCount(sqlTimestamp)
	if err != nil {
		log.Println("[AdminCountOfReceivedMsgsPage_ActionOnInit] CheckMessagesCount_err ==>", err)
		this.setErrorResp("Server error.")
		return
	}
	this.msgCount = msgCount

	this.setErrorResp("")
}

func (this *AdminCountOfReceivedMsgsPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.AdminCountOfReceivedMsgsPageKeyboard
}

var _ models.IPageWithKeyboard = (*AdminCountOfReceivedMsgsPage)(nil)
var _ models.IPageWithActionOnInit = (*AdminCountOfReceivedMsgsPage)(nil)
