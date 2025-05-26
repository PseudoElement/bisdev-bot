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
}

func NewAdminListOfMessagesPage(db models.IDatabase, adminQueryBuilder *query_builder.AdminQueryBuilder) *AdminListOfMessagesPage {
	return &AdminListOfMessagesPage{
		Page: NewPage(db, adminQueryBuilder),
	}
}

func (this *AdminListOfMessagesPage) Name() string {
	return consts.ADMIN_LIST_OF_MESSAGES_PAGE
}

func (this *AdminListOfMessagesPage) RespText(update tgbotapi.Update) string {
	query := this.adminQueryBuilder.GetQueryMsg(this.UserName(update))
	messages, err := this.db.Tables().Messages.GetMesages(query)
	if err != nil {
		log.Println("[AdminListOfMessagesPage_RespText] err in GetMesages: ", err)
		return "Server error."
	}
	if len(messages) == 0 {
		return "No new messages."
	}

	str := bytes.NewBufferString("Here is the list of messages:\n")
	for _, msg := range messages {
		row := fmt.Sprintf("User: %s.\nMessage:\n %v\n\n", msg.UserName, msg.Text)
		str.WriteString(row)
	}

	return str.String()
}

func (this *AdminListOfMessagesPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.AdminListOfLinksPageKeyboard
}

func (this *AdminListOfMessagesPage) NextPage(update tgbotapi.Update, isAdmin bool) models.IPage {
	return NewAdminStartPage(this.db, this.adminQueryBuilder)
}

var _ models.IPageWithKeyboard = (*AdminListOfMessagesPage)(nil)
