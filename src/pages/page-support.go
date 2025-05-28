package pages

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/utils"
)

type SupportPage struct {
	*Page
}

func NewSupportPage(db models.IDatabase, bot *tgbotapi.BotAPI, adminQueryBuilder *query_builder.AdminQueryBuilder) *SupportPage {
	p := &SupportPage{
		Page: NewPage(db, bot, adminQueryBuilder),
	}
	p.setCurrenPage(p)

	return p
}

func (this *SupportPage) Name() string {
	return consts.SUPPORT_PAGE
}

func (this *SupportPage) RespText(update tgbotapi.Update) string {
	return `Sorry to hear you're having trouble 😔. Let me help.
Can you please provide:
- A short description of the issue
- Tx Hash (if relevant)
- Network / Tokens involved
- Device & browser (if on web)
- Screenshot (no more than 1 image per request)

🔧 For faster help, feel free to head to our support Telegram: https://t.me/eobuhow.
Or describe your problem here — I’ll log this and escalate it to our tech support team.`
}

func (this *SupportPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.SupportPageKeyboard
}

func (this *SupportPage) ActionOnDestroy(update tgbotapi.Update) {
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
			log.Println("[SupportPage_ActionOnDestroy] Document_ReadUploadedFile_err ==>", err)
		}

		dbMsg.ImageBlob = buf
	}
	if update.Message.Photo != nil {
		photoSizes := update.Message.Photo
		fileId := photoSizes[len(photoSizes)-1].FileID
		buf, err := utils.ReadUploadedFile(this.bot, fileId)
		if err != nil {
			log.Println("[SupportPage_ActionOnDestroy] ReadUploadedFile_err ==>", err)
		}

		dbMsg.ImageBlob = buf
	}

	err := this.db.Tables().Messages.AddMessage(dbMsg)
	if err != nil {
		log.Println("[SupportPage_ActionOnDestroy] AddMessage err ==> ", err)
		this.setErrorResp("Error on server side trying to save your message. Try to contact support directly: https://t.me/eobuhow.")
	} else {
		this.setErrorResp("")
	}
}

var _ models.IPageWithKeyboard = (*SupportPage)(nil)
var _ models.IPageWithActionOnDestroy = (*SupportPage)(nil)
