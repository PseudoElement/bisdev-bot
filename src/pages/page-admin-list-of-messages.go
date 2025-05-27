package pages

import (
	"bytes"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type AdminListOfMessagesPage struct {
	*Page
	respText string
	messages []models.DB_UserMessage
}

func NewAdminListOfMessagesPage(
	db models.IDatabase,
	bot *tgbotapi.BotAPI,
	adminQueryBuilder *query_builder.AdminQueryBuilder,
) *AdminListOfMessagesPage {
	p := &AdminListOfMessagesPage{
		Page:     NewPage(db, bot, adminQueryBuilder),
		respText: "",
		messages: make([]models.DB_UserMessage, 0, 5),
	}
	p.setCurrenPage(p)

	return p
}

func (this *AdminListOfMessagesPage) Name() string {
	return consts.ADMIN_LIST_OF_MESSAGES_PAGE
}

func (this *AdminListOfMessagesPage) RespText(update tgbotapi.Update) string {
	if len(this.messages) == 0 {
		return "No new messages."
	}

	str := bytes.NewBufferString("Here is the list of messages:\n\n")
	for _, msg := range this.messages {
		row := fmt.Sprintf("User: %s.\nMessage:\n %v\n\n", msg.UserName, msg.Text)
		str.WriteString(row)
	}

	return str.String()
}

func (this *AdminListOfMessagesPage) ActionOnInit(update tgbotapi.Update) {
	query := this.adminQueryBuilder.GetQueryMsg(this.UserName(update))
	messages, err := this.db.Tables().Messages.GetMessages(query)
	this.messages = messages

	if err != nil {
		log.Println("[AdminListOfMessagesPage_ActionOnInit] err in GetMessages: ", err)
		this.setErrorResp("Server error.")
		return
	}

	this.setErrorResp("")
}

func (this *AdminListOfMessagesPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.AdminListOfLinksPageKeyboard
}

func (this *AdminListOfMessagesPage) NextPage(update tgbotapi.Update, isAdmin bool) models.IPage {
	if this.errResp != "" {
		return this
	}
	return NewAdminStartPage(this.db, this.bot, this.adminQueryBuilder)
}

var _ models.IPageWithKeyboard = (*AdminListOfMessagesPage)(nil)
var _ models.IPageWithActionOnInit = (*AdminListOfMessagesPage)(nil)
