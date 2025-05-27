package pages

import (
	"bytes"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type AdminInputUserName struct {
	*Page
	userNames models.DB_UserNames
}

func NewAdminInputUserName(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *AdminInputUserName {
	p := &AdminInputUserName{
		Page: NewPage(db, bot, adminQueryBuilder),
		userNames: models.DB_UserNames{
			NotRead:     make([]string, 0, 10),
			AlreadyRead: make([]string, 0, 10),
		},
	}
	p.setCurrenPage(p)

	return p
}

func (this *AdminInputUserName) Name() string {
	return consts.ADMIN_INPUT_USER_NAME_PAGE
}

func (this *AdminInputUserName) RespText(update tgbotapi.Update) string {
	if this.errResp != "" {
		return this.errResp
	}

	if len(this.userNames.NotRead) > 0 || len(this.userNames.AlreadyRead) > 0 {
		str := bytes.NewBufferString("Input username from the list below:\n\n")

		if len(this.userNames.NotRead) > 0 {
			str.WriteString("New messages from:\n\t")
			str.WriteString(strings.Join(this.userNames.NotRead, ", "))
			str.WriteString(".\n\n")
		}

		str.WriteString("Already read messages from:\n\t")
		str.WriteString(strings.Join(this.userNames.AlreadyRead, ", "))
		str.WriteString(".")

		return str.String()
	} else {
		return "No messages yet."
	}
}

func (this *AdminInputUserName) ActionOnInit(update tgbotapi.Update) {
	userNames, err := this.db.Tables().Messages.GetUserNames()
	if err != nil {
		log.Println("[AdminInputUserName_ActionOnInit] GetUserNames_err ==>", err)
		this.setErrorResp("Server error.")
		return
	}
	this.userNames = userNames

	this.setErrorResp("")
}

// func (this *AdminInputUserName) ActionOnDestroy(update tgbotapi.Update) {
// 	if update.Message == nil {
// 		return
// 	}

// 	this.setErrorResp("")
// 	this.adminQueryBuilder.SetUserNameQueryMsgsOfUsers(
// 		this.UserName(update),
// 		this.TextFromClient(update),
// 	)
// }

func (this *AdminInputUserName) NextPage(update tgbotapi.Update, isAdmin bool) models.IPage {
	if this.errResp != "" {
		return this
	} else {
		return NewAdminListOfSingleUserMsgsPage(this.db, this.bot, this.adminQueryBuilder)
	}
}

// var _ models.IPageWithActionOnDestroy = (*AdminInputUserName)(nil)
var _ models.IPageWithActionOnInit = (*AdminInputUserName)(nil)
