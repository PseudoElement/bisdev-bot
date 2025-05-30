package pages

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

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
	if utf8.RuneCountInString(this.TextFromClient(update)) > 500 {
		this.setErrorResp("Too large message.")
		return
	}

	dbMsg := models.UserMsgFromClient{
		UserName:  this.UserName(update),
		UserId:    this.UserID(update),
		Initials:  fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName),
		Text:      this.TextFromClient(update),
		CreatedAt: utils.GetSqlTimestampByMinutesUTC(0, false),
	}

	doc := update.Message.Document
	if doc != nil {
		if this.isSupportedMimetype(doc) {
			if doc.FileSize > consts.MB_5 {
				this.setErrorResp("Too large file. Max size is 5mb.")
				return
			}

			fileId := update.Message.Document.FileID
			buf, err := utils.ReadUploadedFile(this.bot, fileId)
			if err != nil {
				log.Printf("[%s_ActionOnDestroy] Document_ReadUploadedFile_err ==> %v\n", this.CurrPage().Name(), err)
			}

			dbMsg.Blob = buf
			dbMsg.BlobType = utils.MimeTypeToSqlBlobType(doc.MimeType)
		} else {
			log.Println("mimetype_error")
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

		dbMsg.Blob = buf
		// @TODO check if need handle different image extensions
		dbMsg.BlobType = consts.FILE_TYPES.Png
	}

	this.setErrorResp("")

	go func() {
		err := this.injector.Db.Tables().Messages.AddMessage(dbMsg)
		if err != nil {
			log.Println("[IssueDescriptionPage_ActionOnDestroy] Messages_AddMessage err ==> ", err)
		}
		this.injector.Notifier.NotifyAdminsOnNewMsg(dbMsg)

		err = this.injector.Db.Tables().MessagesCount.AddMessage(dbMsg)
		if err != nil {
			log.Println("[IssueDescriptionPage_ActionOnDestroy] MessagesCount_AddMessage err ==> ", err)
		}
	}()
}

func (this *AbstrUserInputPage) isSupportedMimetype(doc *tgbotapi.Document) bool {
	for _, t := range consts.SUPPORTED_MIME_TYPES {
		if strings.ToLower(t) == strings.ToLower(doc.MimeType) {
			return true
		}
	}
	return false
}
