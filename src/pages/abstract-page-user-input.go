package pages

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/utils"
)

type AbstrUserInputPage struct {
	*Page
}

func NewAbstrUserInputPage(page *Page) *AbstrUserInputPage {
	p := &AbstrUserInputPage{
		Page: page,
	}

	return p
}

func (this *AbstrUserInputPage) ActionOnDestroy(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	if len(this.TextFromClient(update)) > 500 {
		this.setErrorResp("Too long message. Max length is 500 chars")
		return
	}

	dbMsg := models.JsonMsgFromClient{
		UserName: this.UserName(update),
		Initials: fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName),
		Text:     this.TextFromClient(update),
	}

	doc := update.Message.Document
	if doc != nil && (doc.MimeType == "image/jpeg" || doc.MimeType == "image/png") {
		fileId := update.Message.Document.FileID
		buf, err := utils.ReadUploadedFile(this.bot, fileId)
		if err != nil {
			log.Printf("[%s_ActionOnDestroy] Document_ReadUploadedFile_err ==> %v\n", this.CurrPage().Name(), err)
		}

		dbMsg.ImageBlob = buf
	}
	if update.Message.Photo != nil {
		photoSizes := update.Message.Photo
		fileId := photoSizes[len(photoSizes)-1].FileID
		buf, err := utils.ReadUploadedFile(this.bot, fileId)
		if err != nil {
			log.Printf("[%s_ActionOnDestroy] ReadUploadedFile_err ==> %v\n", this.CurrPage().Name(), err)
		}

		dbMsg.ImageBlob = buf
	}

	go func() {
		err := this.db.Tables().Messages.AddMessage(dbMsg)
		if err != nil {
			log.Println("[IssueDescriptionPage_ActionOnDestroy] AddMessage err ==> ", err)
		}
	}()
}
