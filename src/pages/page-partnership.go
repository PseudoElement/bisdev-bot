package pages

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/utils"
)

type PartnershipPage struct {
	*Page
}

func NewPartnershipPage(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *PartnershipPage {
	p := &PartnershipPage{
		Page: NewPage(db, bot, adminQueryBuilder),
	}
	p.setCurrenPage(p)

	return p
}

func (this *PartnershipPage) Name() string {
	return consts.PARTNERSHIP_PAGE
}

func (this *PartnershipPage) RespText(update tgbotapi.Update) string {
	if this.errResp != "" {
		return this.errResp
	}

	return `Awesome!🙌 Let's explore a potential collaboration.

Can you share the following:
- Project name
- Website
- Your role
- Your main goal with us? (integration / liquidity aggregation / mutual routing / co-marketing / other)

Once you're done, I’ll share this with our BD team and we’ll follow up fast.  `
}

func (this *PartnershipPage) ActionOnDestroy(update tgbotapi.Update) {
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
			log.Println("[PartnershipPage_ActionOnDestroy] Document_ReadUploadedFile_err ==>", err)
		}

		dbMsg.ImageBlob = buf
	}
	if update.Message.Photo != nil {
		photoSizes := update.Message.Photo
		fileId := photoSizes[len(photoSizes)-1].FileID
		buf, err := utils.ReadUploadedFile(this.bot, fileId)
		if err != nil {
			log.Println("[PartnershipPage_ActionOnDestroy] ReadUploadedFile_err ==>", err)
		}

		dbMsg.ImageBlob = buf
	}

	err := this.db.Tables().Messages.AddMessage(dbMsg)
	if err != nil {
		log.Println("[PartnershipPage_ActionOnDestroy] AddMessage err ==> ", err)
		this.setErrorResp("Error on server side trying to save your message. Try to contact support directly: https://t.me/eobuhow.")
	} else {
		this.setErrorResp("")
	}
}

var _ models.IPage = (*PartnershipPage)(nil)
var _ models.IPageWithActionOnDestroy = (*PartnershipPage)(nil)
