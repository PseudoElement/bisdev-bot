package pages

import (
	"bytes"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/utils"
)

type AdminListOfMessagesPage struct {
	*Page
	messages []models.DB_UserMessage
}

func NewAdminListOfMessagesPage(injector *injector.AppInjector) *AdminListOfMessagesPage {
	p := &AdminListOfMessagesPage{
		Page:     NewPage(injector),
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
		return "No messages found."
	}

	str := bytes.NewBufferString("Here is the list of messages:\n\n")
	for idx, msg := range this.messages {
		row := fmt.Sprintf("%d. User: %s\n Initials: %s\n Creation time(Moscow time): %v\n Message: %v\n\n",
			idx+1,
			msg.UserName,
			msg.Initials,
			utils.ConvertUTCToMoscowTime(msg.CreatedAt),
			msg.Text,
		)
		str.WriteString(row)
	}

	return str.String()
}

func (this *AdminListOfMessagesPage) ActionOnInit(update tgbotapi.Update) {
	query := this.injector.AdminQueryBuilder.GetQueryMsg(this.UserName(update))
	messages, err := this.injector.Db.Tables().Messages.GetMessages(query)
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
	return NewAdminStartPage(this.injector)
}

var _ models.IPageWithKeyboard = (*AdminListOfMessagesPage)(nil)
var _ models.IPageWithActionOnInit = (*AdminListOfMessagesPage)(nil)
