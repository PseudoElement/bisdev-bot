package pages

import (
	"bytes"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/utils"
)

type AdminInputUserNamePage struct {
	*Page
	userNames   models.DB_UserNames
	commandName string
}

func NewAdminInputUserNamePage(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *AdminInputUserNamePage {
	p := &AdminInputUserNamePage{
		Page: NewPage(db, bot, adminQueryBuilder),
		userNames: models.DB_UserNames{
			NotRead:     make([]string, 0, 10),
			AlreadyRead: make([]string, 0, 10),
		},
		commandName: "",
	}
	p.setCurrenPage(p)

	return p
}

func (this *AdminInputUserNamePage) Name() string {
	return consts.ADMIN_INPUT_USER_NAME_PAGE
}

func (this *AdminInputUserNamePage) AllowedOnlyMessages() bool {
	return true
}

func (this *AdminInputUserNamePage) RespText(update tgbotapi.Update) string {
	if this.errResp != "" {
		return this.errResp
	}
	if this.warningResp != "" {
		return this.warningResp
	}

	if len(this.userNames.NotRead) > 0 || len(this.userNames.AlreadyRead) > 0 {
		if this.commandName == consts.SHOW_MESSAGES_OF_SPECIFIC_USER {
			return this.respForShowMessages(update)
		} else if this.commandName == consts.BLOCK_USER {
			return this.respForBlockUser(update)
		} else {
			return this.respForDeletion(update)
		}
	} else {
		return "No messages yet."
	}
}

func (this *AdminInputUserNamePage) ActionOnDestroy(update tgbotapi.Update) {
	userNameFromInput := this.TextFromClient(update)
	allUserNames := append(this.userNames.AlreadyRead, this.userNames.NotRead...)

	for _, name := range allUserNames {
		if strings.ToLower(name) == strings.ToLower(userNameFromInput) {
			if this.commandName == consts.DELETE_MESSAGES_OF_USER {
				go this.onDestroyDeleteMsgsOfUser(update)
			}
			if this.commandName == consts.BLOCK_USER {
				// not in goroutine cause need check if user is admin
				this.onDestroyBlockUser(update)
			}
			this.setErrorResp("")
			return
		}
	}

	this.setErrorResp("User with name " + userNameFromInput + " doesn't exist.")
	log.Println("[AdminInputUserNamePage_ActionOnDestroy] err =>", this.errResp)
}

func (this *AdminInputUserNamePage) ActionOnInit(update tgbotapi.Update) {
	userNames, err := this.db.Tables().Messages.GetUserNames()
	if err != nil {
		log.Println("[AdminInputUserNamePage_ActionOnInit] GetUserNames_err ==>", err)
		this.setErrorResp("Server error.")
		return
	}
	this.userNames = userNames
	this.commandName = this.TextFromClient(update)
}

func (this *AdminInputUserNamePage) NextPage(update tgbotapi.Update, isAdmin bool) models.IPage {
	if this.TextFromClient(update) == consts.BACK_TO_START {
		return NewAdminStartPage(this.db, this.bot, this.adminQueryBuilder)
	} else if this.errResp != "" || this.warningResp != "" {
		return this
	} else {
		if this.commandName == consts.DELETE_MESSAGES_OF_USER {
			return NewAdminInfoAfterDeletionPage(this.db, this.bot, this.adminQueryBuilder)
		} else if this.commandName == consts.BLOCK_USER {
			return NewNotificationAfterBlockUserPage(this.db, this.bot, this.adminQueryBuilder)
		} else {
			// consts.SHOW_MESSAGES_OF_SPECIFIC_USER
			return NewAdminListOfSingleUserMsgsPage(this.db, this.bot, this.adminQueryBuilder)
		}
	}
}

func (this *AdminInputUserNamePage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.BackToStartKeyBoard
}

func (this *AdminInputUserNamePage) respForDeletion(update tgbotapi.Update) string {
	allUserNames := append(this.userNames.AlreadyRead, this.userNames.NotRead...)
	uniqueNames := utils.FilterUnique(allUserNames)
	str := bytes.NewBufferString("Input username you want to delete from the list below:\n\n")
	str.WriteString(strings.Join(uniqueNames, ", "))

	return str.String()
}

func (this *AdminInputUserNamePage) respForShowMessages(update tgbotapi.Update) string {
	str := bytes.NewBufferString("Input username you want to check from the list below:\n\n")
	if len(this.userNames.NotRead) > 0 {
		str.WriteString("New messages from:\n")
		str.WriteString(strings.Join(this.userNames.NotRead, ", "))
		str.WriteString("\n\n")
	}

	if len(this.userNames.AlreadyRead) > 0 {
		str.WriteString("Already read messages from:\n")
		str.WriteString(strings.Join(this.userNames.AlreadyRead, ", "))
	}

	return str.String()
}

func (this *AdminInputUserNamePage) respForBlockUser(update tgbotapi.Update) string {
	allUserNames := append(this.userNames.AlreadyRead, this.userNames.NotRead...)
	uniqueNames := utils.FilterUnique(allUserNames)
	str := bytes.NewBufferString("Input username you want to block from the list below:\n\n")
	str.WriteString(strings.Join(uniqueNames, ", "))

	return str.String()
}

// here userName 100% existing
func (this *AdminInputUserNamePage) onDestroyBlockUser(update tgbotapi.Update) {
	userName := this.TextFromClient(update)
	isAdmin := this.IsUserAdmin(userName)
	if isAdmin {
		this.setErrorResp(userName + "is admin. You're not allowed to block another admin.")
		return
	}

	go func() {
		err := this.db.Tables().BlockedUsers.BlockUser(userName)
		if err != nil {
			log.Println("[AdminInputUserNamePage_onDestroyBlockUser] BlockUser_err ==>", err)
		}
	}()

	this.setErrorResp("")
}

// here userName 100% existing
func (this *AdminInputUserNamePage) onDestroyDeleteMsgsOfUser(update tgbotapi.Update) {
	userName := this.TextFromClient(update)
	err := this.db.Tables().Messages.DeleteMessagesByUserName(userName)
	if err != nil {
		log.Println("[AdminInputUserNamePage_onDestroyDeleteMsgsOfUser] Delete_err ==>", err)
	}
}

var _ models.IPageWithActionOnDestroy = (*AdminInputUserNamePage)(nil)
var _ models.IPageWithActionOnInit = (*AdminInputUserNamePage)(nil)
var _ models.IPageWithKeyboard = (*AdminInputUserNamePage)(nil)
