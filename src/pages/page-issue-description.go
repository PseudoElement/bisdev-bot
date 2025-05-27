package pages

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/utils"
)

type IssueDescriptionPage struct {
	*Page
	errResp string
}

func NewIssueDescriptionPage(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *IssueDescriptionPage {
	p := &IssueDescriptionPage{
		Page:    NewPage(db, bot, adminQueryBuilder),
		errResp: "",
	}
	p.setCurrenPage(p)

	return p
}

func (this *IssueDescriptionPage) Name() string {
	return consts.DESCRIBE_ISSUE
}

func (this *IssueDescriptionPage) RespText(update tgbotapi.Update) string {
	return `Can you please provide:
- A short description of the issue
- Tx Hash (if relevant)
- Network / Tokens involved
- Device & browser (if on web)
- Screenshot (if possible)`
}

func (this *IssueDescriptionPage) ActionOnDestroy(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	dbMsg := models.JsonMsgFromClient{
		UserName: this.UserName(update),
		Text:     this.TextFromClient(update),
	}
	if update.Message.Photo != nil {
		photoSizes := update.Message.Photo
		fileId := photoSizes[len(photoSizes)-1].FileID
		buf, err := utils.ReadUploadedFile(this.bot, fileId)
		if err != nil {
			log.Println("[IssueDescriptionPage_Action] ReadUploadedFile_err ==>", err)
		}

		dbMsg.ImageBlob = buf
	}

	err := this.db.Tables().Messages.AddMessage(dbMsg)
	if err != nil {
		log.Println("[IssueDescriptionPage_Action] AddMessage err ==> ", err)
		this.setErrorResp("Error on server side trying to save your message. Try to contact support directly: https://t.me/eobuhow.")
	} else {
		this.setErrorResp("")
	}
}

var _ models.IPageWithActionOnDestroy = (*IssueDescriptionPage)(nil)
