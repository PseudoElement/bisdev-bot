package pages

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
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

func (this *AbstrUserInputPage) AllowedOnlyMessages() bool {
	return true
}

func (this *AbstrUserInputPage) ActionOnDestroy(update tgbotapi.Update) {
	if update.CallbackData() == consts.BACK_TO_START {
		this.setErrorResp("")
		return
	}
	if update.Message == nil {
		return
	}
	if len(this.TextFromClient(update)) > 500 {
		this.setErrorResp("Too long message. Max length is 500 chars")
		return
	}

	dbMsg := models.JsonMsgFromClient{
		UserName:  this.UserName(update),
		Initials:  fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName),
		Text:      this.TextFromClient(update),
		CreatedAt: utils.GetSqlTimestampByMinutesUTC(0, false),
	}

	doc := update.Message.Document
	if doc != nil {
		if doc.MimeType == "image/jpeg" || doc.MimeType == "image/png" {
			if doc.FileSize > consts.MB_5 {
				this.setErrorResp("Too large file. Max size is 5mb.")
				return
			}

			fileId := update.Message.Document.FileID
			buf, err := utils.ReadUploadedFile(this.bot, fileId)
			if err != nil {
				log.Printf("[%s_ActionOnDestroy] Document_ReadUploadedFile_err ==> %v\n", this.CurrPage().Name(), err)
			}

			dbMsg.ImageBlob = buf
		} else {
			this.setErrorResp(doc.MimeType + " file format is not supported.")
			return
		}

	}
	if update.Message.Photo != nil {
		photoSizes := update.Message.Photo
		bestQualityPic := photoSizes[len(photoSizes)-1]
		if bestQualityPic.FileSize > consts.MB_5 {
			this.setErrorResp("Too large picture. Max size is 5mb.")
			return
		}

		fileId := bestQualityPic.FileID
		buf, err := utils.ReadUploadedFile(this.bot, fileId)
		if err != nil {
			log.Printf("[%s_ActionOnDestroy] ReadUploadedFile_err ==> %v\n", this.CurrPage().Name(), err)
		}

		dbMsg.ImageBlob = buf
	}

	this.setErrorResp("")

	go func() {
		err := this.db.Tables().Messages.AddMessage(dbMsg)
		if err != nil {
			log.Println("[IssueDescriptionPage_ActionOnDestroy] Messages_AddMessage err ==> ", err)
		}
		err = this.db.Tables().MessagesCount.AddMessage(dbMsg)
		if err != nil {
			log.Println("[IssueDescriptionPage_ActionOnDestroy] MessagesCount_AddMessage err ==> ", err)
		}
	}()
}
