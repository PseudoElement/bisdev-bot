package bd_bot

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/db"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type BuisdevBot struct {
	bot               *tgbotapi.BotAPI
	db                *db.SqliteDB
	isProd            bool
	page              models.IPage
	lastCommand       string
	admins            []string
	adminQueryBuilder *query_builder.AdminQueryBuilder
}

func NewBuisdevBot() *BuisdevBot {
	token, ok := os.LookupEnv("BOT_API_KEY")
	if !ok {
		panic("BOT_API_KEY variable not provided in .env file.")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	isProd, err := strconv.ParseBool(os.Getenv("IS_PROD"))
	if err != nil {
		panic("IS_PROD variable supposed to be true of false.")
	}

	adminsString, ok := os.LookupEnv("ADMINS")
	if !ok {
		adminsString = ""
	}
	admins := strings.Split(adminsString, " ")

	bot.Debug = !isProd

	db := db.NewSqliteDB()
	adminQueryBuilder := query_builder.NewAdminQueryBuilder()

	b := &BuisdevBot{
		bot:               bot,
		isProd:            isProd,
		page:              pages.NewStartPage(db, bot, adminQueryBuilder),
		admins:            admins,
		lastCommand:       "",
		db:                db,
		adminQueryBuilder: adminQueryBuilder,
	}

	return b
}

func (this *BuisdevBot) Listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range this.bot.GetUpdatesChan(u) {
		if update.Message != nil {
			fmt.Printf("[%s][%s %s] userId - %v, text - %s,  caption - %s, command - %s, photos_count - %d.\n",
				update.Message.From.UserName,
				update.Message.From.FirstName,
				update.Message.From.LastName,
				update.Message.From.ID,
				update.Message.Text,
				update.Message.Caption,
				update.Message.Command(),
				len(update.Message.Photo),
			)
		}

		if update.Message != nil {
			// log.Printf("animation ==> %+v\n", update.Message.Animation)
			// log.Printf("callback_data ==> %+v\n", update.CallbackData())
			// log.Printf("mediagroup_id ==> %+v\n", update.Message.MediaGroupID)
			// log.Printf("photo ==> %+v\n", update.Message.Photo[0])

			// pwd, _ := os.Getwd()
			// path := pwd + "/src/bot/beach.jpg"

			// path2 := pwd + "/src/bot/rbc-util-icon.png"
			// path3 := pwd + "/src/bot/usdc-util-icon.png"
			// photoFromServer := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath(path))
			// photoFromServer2 := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath(path2))
			// photoFromServer3 := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath(path3))

			// photoFromClient := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(update.Message.Photo[0].FileID))

			// filePath := pwd + "/src/bot/test_results.json"
			// file, err := os.Open(filePath)
			// log.Println("ReadFile_err ==> ", err)

			// fileReader := tgbotapi.FileReader{Name: "results_new.json", Reader: file}
			// fileFromServer := tgbotapi.NewInputMediaDocument(fileReader)
			// fileMediaGroup := tgbotapi.NewMediaGroup(update.Message.Chat.ID, []interface{}{
			// 	fileFromServer,
			// })
			// fileFromClient := tgbotapi.NewInputMediaDocument(tgbotapi.FileID(update.Message.Document.FileID))

			// buf, err := utils.ReadUploadedFile(this.bot, update.Message.Document.FileID)

			// log.Println("Upload_err ==> ", err)
			// log.Println("buf ==> ", buf)

			// go func() {
			// 	mediaGroup := tgbotapi.NewMediaGroup(update.Message.Chat.ID, []interface{}{
			// 		// photoFromServer,
			// 		// photoFromServer2,
			// 		// photoFromServer3,
			// 		// fileFromServer,
			// 		// photoFromClient,
			// 		// fileFromServer,
			// 		// fileFromClient,
			// 		// tgbotapi.NewInputMediaDocument(tgbotapi.FileBytes{Name: "parsed_file.doc", Bytes: buf}),
			// 	})

			// 	this.bot.Send(mediaGroup)
			// 	log.Println("after media group")

			// 	// this.bot.Send(fileMediaGroup)
			// 	// log.Println("after file media group")

			// 	// this.bot.GetFile(tgbotapi.FileConfig{FileID: update.Message.Animation.FileID})
			// }()

			// continue
		}

		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			userId := update.Message.From.ID

			switch update.Message.Text {
			case "/start":
				if this.isAdmin(userId) {
					this.page = pages.NewAdminStartPage(this.db, this.bot, this.adminQueryBuilder)
				} else {
					this.page = pages.NewStartPage(this.db, this.bot, this.adminQueryBuilder)
				}
				msg.Text = this.page.RespText(update)
				msg.ReplyMarkup = this.page.(models.IPageWithKeyboard).Keyboard()

				go this.bot.Send(msg)
			default:
				if update.Message.Text == "" && update.Message.Caption == "" {
					msg.Text = "Message wihout text is not allowed to save."

					go this.bot.Send(msg)
				} else if this.page.AllowedOnlyCommands() {
					msg.Text = "Select one option from the list."
					msg.ReplyMarkup = this.page.(models.IPageWithKeyboard).Keyboard()

					go this.bot.Send(msg)
				} else {
					pageWithActionOnDestr, withOnDestroy := this.page.(models.IPageWithActionOnDestroy)
					if withOnDestroy {
						pageWithActionOnDestr.ActionOnDestroy(update)
					}

					nextPage := this.page.NextPage(update, this.isAdmin(userId))
					nextPageWithActionOnInit, withOnInit := nextPage.(models.IPageWithActionOnInit)
					nextPageWithKeyboard, withKeyboard := nextPage.(models.IPageWithKeyboard)
					nextPageWithPhotos, withPhotosInResp := nextPage.(models.IPageWithPhotos)

					if withOnInit {
						nextPageWithActionOnInit.ActionOnInit(update)
					}
					if withKeyboard {
						msg.ReplyMarkup = nextPageWithKeyboard.Keyboard()
					}
					msg.Text = nextPage.RespText(update)

					go func() {
						this.bot.Send(msg)
						if withPhotosInResp && nextPageWithPhotos.HasPhotos() {
							this.bot.SendMediaGroup(nextPageWithPhotos.PhotosResp(update))
						}
					}()

					this.page = nextPage
				}
			}
		} else if update.CallbackQuery != nil && update.CallbackData() != this.lastCommand {
			fmt.Printf("[%s][%s %s] userId - %v, Data - %s\n",
				update.CallbackQuery.From.UserName,
				update.CallbackQuery.From.FirstName,
				update.CallbackQuery.From.LastName,
				update.CallbackQuery.From.ID,
				update.CallbackData(),
			)
			this.lastCommand = update.CallbackData()
			userId := update.CallbackQuery.From.ID

			pageWithActionOnDestr, ok := this.page.(models.IPageWithActionOnDestroy)
			if ok {
				pageWithActionOnDestr.ActionOnDestroy(update)
			}

			nextPage := this.page.NextPage(update, this.isAdmin(userId))
			nextPageWithActionOnInit, ok := nextPage.(models.IPageWithActionOnInit)
			if ok {
				nextPageWithActionOnInit.ActionOnInit(update)
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, nextPage.RespText(update))

			nextPageWithKB, ok := nextPage.(models.IPageWithKeyboard)
			if ok {
				msg.ReplyMarkup = nextPageWithKB.Keyboard()
			}

			this.page = nextPage

			go this.bot.Send(msg)
		}
	}
}

func (this *BuisdevBot) isAdmin(userId int64) bool {
	for _, adminId := range this.admins {
		if strings.ToLower(strconv.Itoa(int(userId))) == strings.ToLower(adminId) {
			return true
		}
	}

	return false
}
