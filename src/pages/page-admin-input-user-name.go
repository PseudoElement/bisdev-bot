package pages

import (
	"bytes"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/utils"
)

type AdminInputUserNamePage struct {
	*Page
	userNames   models.DB_UserNames
	commandName string
}

func NewAdminInputUserNamePage(injector *injector.AppInjector) *AdminInputUserNamePage {
	p := &AdminInputUserNamePage{
		Page: NewPage(injector),
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
		} else if this.commandName == consts.UNBLOCK_USER {
			return this.respForUnblockUser(update)
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
			if this.commandName == consts.UNBLOCK_USER {
				go this.onDestroyUnblockUser(update)
			}
			if this.commandName == consts.BLOCK_USER {
				// not in goroutine cause need check if user is admin
				this.onDestroyBlockUser(update)
			}

			return
		}
	}

	this.setErrorResp("User with name " + userNameFromInput + " doesn't exist.")
	log.Println("[AdminInputUserNamePage_ActionOnDestroy] err =>", this.errResp)
}

func (this *AdminInputUserNamePage) ActionOnInit(update tgbotapi.Update) {
	if len(this.userNames.AlreadyRead) == 0 && len(this.userNames.NotRead) == 0 {
		if update.CallbackData() != "" {
			this.commandName = this.TextFromClient(update)
		}

		if this.commandName == consts.UNBLOCK_USER {
			blockedUsers := this.injector.Store.GetBlockedUsers()
			if len(blockedUsers) < 1 {
				this.setWarningResp("Blacklist is empty. There is no one to unblock.")
				return
			}

			for _, blockedUser := range blockedUsers {
				this.userNames.AlreadyRead = append(this.userNames.AlreadyRead, blockedUser.UserName)
			}
		} else {
			userNames, err := this.injector.Db.Tables().Messages.GetUserNames()
			if err != nil {
				log.Println("[AdminInputUserNamePage_ActionOnInit] GetUserNames_err ==>", err)
				this.setErrorResp("Server error.")
				return
			}
			this.userNames = userNames
		}

	}

}

func (this *AdminInputUserNamePage) NextPage(update tgbotapi.Update, isAdmin bool) models.IPage {
	log.Println("[AdminInputUserNamePage_NextPage] commandName - ", this.commandName, " err - ", this.errResp)
	if this.TextFromClient(update) == consts.BACK_TO_START {
		return NewAdminStartPage(this.injector)
	} else if this.errResp != "" || this.warningResp != "" {
		return this
	} else {
		if this.commandName == consts.DELETE_MESSAGES_OF_USER {
			return NewAdminInfoAfterDeletionPage(this.injector)
		} else if this.commandName == consts.BLOCK_USER {
			return NewNotificationAfterBlockUserPage(this.injector)
		} else if this.commandName == consts.UNBLOCK_USER {
			return NewAdminNotificationAfterUserUnblockPage(this.injector)
		} else {
			// consts.SHOW_MESSAGES_OF_SPECIFIC_USER
			return NewAdminListOfSingleUserMsgsPage(this.injector)
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

func (this *AdminInputUserNamePage) respForUnblockUser(update tgbotapi.Update) string {
	allUserNames := append(this.userNames.AlreadyRead, this.userNames.NotRead...)
	uniqueNames := utils.FilterUnique(allUserNames)
	str := bytes.NewBufferString("Input username you want to unblock from the blacklist below:\n\n")
	str.WriteString(strings.Join(uniqueNames, ", "))

	return str.String()
}

// here userName 100% existing
func (this *AdminInputUserNamePage) onDestroyBlockUser(update tgbotapi.Update) {
	userName := this.TextFromClient(update)
	if this.injector.Store.IsAdminByName(userName) {
		log.Printf("!!!NOTE: %s tried to block admin %s.\n", update.Message.From.UserName, userName)
		this.setErrorResp(userName + " is admin. You're not allowed to block another admin.")
		return
	}

	go func() {
		err := this.injector.Db.Tables().BlockedUsers.BlockUser(userName)
		if err != nil {
			log.Println("[AdminInputUserNamePage_onDestroyBlockUser] BlockUser_err ==>", err)
		}
		this.injector.Notifier.NotifyAdminsOnBlockedUsers(userName, this.UserName(update))
		this.injector.Store.UpdateBlockedUsersList()
	}()

	this.setErrorResp("")
}

// here userName 100% existing
func (this *AdminInputUserNamePage) onDestroyUnblockUser(update tgbotapi.Update) {
	userName := this.TextFromClient(update)
	err := this.injector.Db.Tables().BlockedUsers.UnblockUser(userName)
	if err != nil {
		log.Println("[AdminInputUserNamePage_onDestroyUnblockUser] UnblockUser_err ==>", err)
	}
	this.injector.Notifier.NotifyAdminsOnUnblockedUsers(userName, this.UserName(update))
	this.injector.Store.UpdateBlockedUsersList()

	this.setErrorResp("")
}

// here userName 100% existing
func (this *AdminInputUserNamePage) onDestroyDeleteMsgsOfUser(update tgbotapi.Update) {
	userName := this.TextFromClient(update)
	err := this.injector.Db.Tables().Messages.DeleteMessagesByUserName(userName)
	if err != nil {
		log.Println("[AdminInputUserNamePage_onDestroyDeleteMsgsOfUser] Delete_err ==>", err)
	}

	this.setErrorResp("")
}

var _ models.IPageWithActionOnDestroy = (*AdminInputUserNamePage)(nil)
var _ models.IPageWithActionOnInit = (*AdminInputUserNamePage)(nil)
var _ models.IPageWithKeyboard = (*AdminInputUserNamePage)(nil)
