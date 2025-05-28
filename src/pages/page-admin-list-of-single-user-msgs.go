package pages

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type AdminListOfSingleUserMsgsPage struct {
	*Page
	respText string
	messages []models.DB_UserMessage
}

func NewAdminListOfSingleUserMsgsPage(
	db models.IDatabase,
	bot *tgbotapi.BotAPI,
	adminQueryBuilder *query_builder.AdminQueryBuilder,
) *AdminListOfSingleUserMsgsPage {
	p := &AdminListOfSingleUserMsgsPage{
		Page:     NewPage(db, bot, adminQueryBuilder),
		respText: "",
		messages: make([]models.DB_UserMessage, 0, 5),
	}
	p.setCurrenPage(p)

	return p
}

func (this *AdminListOfSingleUserMsgsPage) Name() string {
	return consts.ADMIN_LIST_OF_SINGLE_USER_MESSAGES_PAGE
}

func (this *AdminListOfSingleUserMsgsPage) HasPhotos() bool {
	for _, msg := range this.messages {
		if msg.ImgBlob != nil && len(msg.ImgBlob) > 0 {
			return true
		}
	}

	return false
}

func (this *AdminListOfSingleUserMsgsPage) RespText(update tgbotapi.Update) string {
	if this.errResp != "" {
		return this.errResp
	}

	if len(this.messages) == 0 {
		return "No messages found."
	}

	str := bytes.NewBufferString("Here is the list of messages:\n")
	for _, msg := range this.messages {
		row := fmt.Sprintf("User: %s.\nMessage:\n %v\n\n", msg.UserName, msg.Text)
		str.WriteString(row)
	}

	return str.String()
}

// @TODO handle more than 10 photos in resp
func (this *AdminListOfSingleUserMsgsPage) PhotosResp(update tgbotapi.Update) tgbotapi.MediaGroupConfig {
	photos := make([]interface{}, 0, len(this.messages))

	idx := 1
	for _, msg := range this.messages {
		// 10 photos max
		if idx > 10 {
			break
		}

		if msg.ImgBlob != nil && len(msg.ImgBlob) > 0 {
			buf := msg.ImgBlob
			fileName := "img_" + strconv.Itoa(idx) + ".png"
			fileBytes := tgbotapi.FileBytes{Name: fileName, Bytes: buf}
			photos = append(photos, tgbotapi.NewInputMediaPhoto(fileBytes))
		}

		idx++
	}

	photoMG := tgbotapi.NewMediaGroup(update.Message.Chat.ID, photos)

	return photoMG
}

func (this *AdminListOfSingleUserMsgsPage) ActionOnInit(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	messages, err := this.db.Tables().Messages.GetMessagesByUserName(this.TextFromClient(update))
	this.messages = messages

	if err != nil {
		log.Println("[AdminListOfSingleUserMsgsPage_ActionOnInit] err in GetMessagesByUserName: ", err)
		this.setErrorResp("Data not found about user " + this.TextFromClient(update) + ".")
		return
	}

	this.setErrorResp("")
}

func (this *AdminListOfSingleUserMsgsPage) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return keyboards.AdminListOfLinksPageKeyboard
}

func (this *AdminListOfSingleUserMsgsPage) NextPage(update tgbotapi.Update, isAdmin bool) models.IPage {
	if this.errResp != "" {
		return this
	}
	return NewAdminStartPage(this.db, this.bot, this.adminQueryBuilder)
}

var _ models.IPageWithKeyboard = (*AdminListOfSingleUserMsgsPage)(nil)
var _ models.IPageWithActionOnInit = (*AdminListOfSingleUserMsgsPage)(nil)
var _ models.IPageWithPhotos = (*AdminListOfSingleUserMsgsPage)(nil)
