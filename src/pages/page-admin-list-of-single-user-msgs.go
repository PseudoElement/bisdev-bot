package pages

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages/keyboards"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/utils"
)

type AdminListOfSingleUserMsgsPage struct {
	*Page
	messages []models.DB_UserMessage
}

func NewAdminListOfSingleUserMsgsPage(
	db models.IDatabase,
	bot *tgbotapi.BotAPI,
	adminQueryBuilder *query_builder.AdminQueryBuilder,
) *AdminListOfSingleUserMsgsPage {
	p := &AdminListOfSingleUserMsgsPage{
		Page:     NewPage(db, bot, adminQueryBuilder),
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
		if this.isImg(msg.BlobType) && msg.Blob != nil && len(msg.Blob) > 0 {
			return true
		}
	}

	return false
}

func (this *AdminListOfSingleUserMsgsPage) HasFiles() bool {
	for _, msg := range this.messages {
		if this.isDoc(msg.BlobType) && msg.Blob != nil && len(msg.Blob) > 0 {
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
	for idx, msg := range this.messages {
		// msg.
		row := fmt.Sprintf("%d. Username: %s\nInitials: %s\nCreation time(Moscow time): %v\nMessage:\n %v\n\n",
			idx+1,
			msg.UserName,
			msg.Initials,
			utils.ConvertUTCToMoscowTime(msg.CreatedAt),
			msg.Text,
		)
		str.WriteString(row)
	}

	return str.String()
}

func (this *AdminListOfSingleUserMsgsPage) FilesResp(update tgbotapi.Update) []tgbotapi.MediaGroupConfig {
	filesChunks := make(
		[][]interface{},
		0,
		int(math.Ceil(float64(len(this.messages)/10))),
	)

	currFilesChunk := 0
	filesChunks = append(filesChunks, make([]interface{}, 0, 10))
	for idx, msg := range this.messages {
		if len(filesChunks[currFilesChunk]) >= 10 {
			currFilesChunk++
			filesChunks = append(filesChunks, make([]interface{}, 0, 10))
		}

		if this.isDoc(msg.BlobType) && msg.Blob != nil && len(msg.Blob) > 0 {
			buf := msg.Blob
			fileName := "file_" + strconv.Itoa(idx+1) + "." + msg.BlobType
			fileBytes := tgbotapi.FileBytes{Name: fileName, Bytes: buf}
			document := tgbotapi.NewInputMediaDocument(fileBytes)

			filesChunks[currFilesChunk] = append(filesChunks[currFilesChunk], document)
		}
	}

	// check if approach by index to filesChunks not buggy
	mediagroups := make([]tgbotapi.MediaGroupConfig, len(filesChunks), len(filesChunks))
	for idx, filesChunk := range filesChunks {
		mg := tgbotapi.NewMediaGroup(update.Message.Chat.ID, filesChunk)
		mediagroups[idx] = mg
	}

	return mediagroups
}

// @TODO handle more than 10 photos in resp
func (this *AdminListOfSingleUserMsgsPage) PhotosResp(update tgbotapi.Update) []tgbotapi.MediaGroupConfig {
	photosChunks := make(
		[][]interface{},
		0,
		int(math.Ceil(float64(len(this.messages)/10))),
	)

	currPhotosChunk := 0
	photosChunks = append(photosChunks, make([]interface{}, 0, 10))
	for idx, msg := range this.messages {
		if len(photosChunks[currPhotosChunk]) >= 10 {
			currPhotosChunk++
			photosChunks = append(photosChunks, make([]interface{}, 0, 10))
		}

		if this.isImg(msg.BlobType) && msg.Blob != nil && len(msg.Blob) > 0 {
			buf := msg.Blob
			fileName := "img_" + strconv.Itoa(idx+1) + "." + msg.BlobType
			fileBytes := tgbotapi.FileBytes{Name: fileName, Bytes: buf}
			document := tgbotapi.NewInputMediaDocument(fileBytes)

			photosChunks[currPhotosChunk] = append(photosChunks[currPhotosChunk], document)
		}
	}

	// check if approach by index to filesChunks not buggy
	mediagroups := make([]tgbotapi.MediaGroupConfig, len(photosChunks), len(photosChunks))
	for idx, filesChunk := range photosChunks {
		mg := tgbotapi.NewMediaGroup(update.Message.Chat.ID, filesChunk)
		mediagroups[idx] = mg
	}

	return mediagroups
}

func (this *AdminListOfSingleUserMsgsPage) ActionOnInit(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	userName := this.TextFromClient(update)

	messages, err := this.db.Tables().Messages.GetMessagesByUserName(userName)
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

func (this *AdminListOfSingleUserMsgsPage) isImg(blobType string) bool {
	for _, t := range consts.IMAGES_FILE_TYPES {
		if t == blobType {
			return true
		}
	}
	return false
}

func (this *AdminListOfSingleUserMsgsPage) isDoc(blobType string) bool {
	for _, t := range consts.DOC_FILE_TYPES {
		if t == blobType {
			return true
		}
	}
	return false
}

var _ models.IPageWithKeyboard = (*AdminListOfSingleUserMsgsPage)(nil)
var _ models.IPageWithActionOnInit = (*AdminListOfSingleUserMsgsPage)(nil)
var _ models.IPageWithPhotos = (*AdminListOfSingleUserMsgsPage)(nil)
var _ models.IPageWithFiles = (*AdminListOfSingleUserMsgsPage)(nil)
