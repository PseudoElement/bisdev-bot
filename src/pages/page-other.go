package pages

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/utils"
)

type OtherPage struct {
	*Page
}

func NewOtherPage(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *OtherPage {
	p := &OtherPage{
		Page: NewPage(db, bot, adminQueryBuilder),
	}
	p.setCurrenPage(p)

	return p
}

func (this *OtherPage) Name() string {
	return consts.OTHER_PAGE
}

func (this *OtherPage) RespText(update tgbotapi.Update) string {
	if this.errResp != "" {
		return this.errResp
	}

	return `No problem! Please describe your request in a few words — I’ll make sure it reaches the right person on our team.
We aim to reply within 24 hours.`
}

func (this *OtherPage) ActionOnDestroy(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	dbMsg := models.JsonMsgFromClient{
		UserName: this.UserName(update),
		Text:     this.TextFromClient(update),
	}
	if update.Message.Document != nil {
		fileId := update.Message.Document.FileID
		buf, err := utils.ReadUploadedFile(this.bot, fileId)
		if err != nil {
			log.Println("[OtherPage_ActionOnDestroy] Document_ReadUploadedFile_err ==>", err)
		}

		dbMsg.ImageBlob = buf
	}
	if update.Message.Photo != nil {
		photoSizes := update.Message.Photo
		fileId := photoSizes[len(photoSizes)-1].FileID
		buf, err := utils.ReadUploadedFile(this.bot, fileId)
		if err != nil {
			log.Println("[OtherPage_ActionOnDestroy] Photo_ReadUploadedFile_err ==>", err)
		}

		dbMsg.ImageBlob = buf
	}

	err := this.db.Tables().Messages.AddMessage(dbMsg)
	if err != nil {
		log.Println("[OtherPage_ActionOnDestroy] AddMessage err ==> ", err)
		this.setErrorResp("Error on server side trying to save your message. Try to contact with support directly: https://t.me/eobuhow.")
	} else {
		this.setErrorResp("")
	}
}

var _ models.IPage = (*OtherPage)(nil)
var _ models.IPageWithActionOnDestroy = (*OtherPage)(nil)
