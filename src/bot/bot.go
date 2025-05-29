package bd_bot

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/db"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
)

type BuisdevBot struct {
	bot               *tgbotapi.BotAPI
	db                *db.SqliteDB
	isProd            bool
	lastCommand       string
	admins            []string
	adminQueryBuilder *query_builder.AdminQueryBuilder
	// key is userId
	pages map[int64]models.IPage
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
		pages:             make(map[int64]models.IPage, 10),
		admins:            admins,
		lastCommand:       "",
		db:                db,
		adminQueryBuilder: adminQueryBuilder,
	}

	return b
}

func (this *BuisdevBot) Listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = math.MaxInt

	for update := range this.bot.GetUpdatesChan(u) {
		var userId int64
		if update.Message != nil {
			userId = update.Message.From.ID
		}
		if update.CallbackQuery != nil {
			userId = update.CallbackQuery.From.ID
		}

		_, ok := this.pages[userId]
		if !ok {
			if this.isAdmin(userId) {
				this.pages[userId] = pages.NewAdminStartPage(this.db, this.bot, this.adminQueryBuilder)
			} else {
				this.pages[userId] = pages.NewStartPage(this.db, this.bot, this.adminQueryBuilder)
			}
		}

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

			this.handleMessageRequest(update)
		} else if update.CallbackQuery != nil && update.CallbackData() != this.lastCommand {
			fmt.Printf("[%s][%s %s] userId - %v, Data - %s\n",
				update.CallbackQuery.From.UserName,
				update.CallbackQuery.From.FirstName,
				update.CallbackQuery.From.LastName,
				update.CallbackQuery.From.ID,
				update.CallbackData(),
			)

			this.lastCommand = update.CallbackData()
			this.handleCallbackRequest(update)
		}
	}
}

func (this *BuisdevBot) handleMessageRequest(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	userId := update.Message.From.ID

	switch update.Message.Text {
	case "/start":
		if this.isAdmin(userId) {
			this.pages[userId] = pages.NewAdminStartPage(this.db, this.bot, this.adminQueryBuilder)
		} else {
			this.pages[userId] = pages.NewStartPage(this.db, this.bot, this.adminQueryBuilder)
		}
		msg.Text = this.pages[userId].RespText(update)
		msg.ReplyMarkup = this.pages[userId].(models.IPageWithKeyboard).Keyboard()

		go this.bot.Send(msg)
	default:
		if update.Message.Text == "" && update.Message.Caption == "" {
			msg.Text = "Message wihout text is not allowed to save."

			go this.bot.Send(msg)
		} else if this.pages[userId].AllowedOnlyCommands() {
			msg.Text = "Select one option from the list."
			msg.ReplyMarkup = this.pages[userId].(models.IPageWithKeyboard).Keyboard()

			go this.bot.Send(msg)
		} else {
			pageWithActionOnDestr, withOnDestroy := this.pages[userId].(models.IPageWithActionOnDestroy)
			if withOnDestroy {
				pageWithActionOnDestr.ActionOnDestroy(update)
			}

			nextPage := this.pages[userId].NextPage(update, this.isAdmin(userId))
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

			this.pages[userId] = nextPage
		}
	}
}

func (this *BuisdevBot) handleCallbackRequest(update tgbotapi.Update) {
	userId := update.CallbackQuery.From.ID

	if this.pages[userId].AllowedOnlyMessages() && update.CallbackData() != consts.BACK_TO_START {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Commands are not allowed for response here.")

		go this.bot.Send(msg)
	} else {
		pageWithActionOnDestr, ok := this.pages[userId].(models.IPageWithActionOnDestroy)
		if ok {
			pageWithActionOnDestr.ActionOnDestroy(update)
		}

		nextPage := this.pages[userId].NextPage(update, this.isAdmin(userId))
		nextPageWithActionOnInit, ok := nextPage.(models.IPageWithActionOnInit)
		if ok {
			nextPageWithActionOnInit.ActionOnInit(update)
		}

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, nextPage.RespText(update))

		nextPageWithKB, ok := nextPage.(models.IPageWithKeyboard)
		if ok {
			msg.ReplyMarkup = nextPageWithKB.Keyboard()
		}

		this.pages[userId] = nextPage

		go this.bot.Send(msg)
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
