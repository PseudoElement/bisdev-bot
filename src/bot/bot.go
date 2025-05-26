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
		page:              pages.NewStartPage(db, adminQueryBuilder),
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
			fmt.Printf("[%s] text - %s, command - %s.\n", update.Message.From.UserName, update.Message.Text, update.Message.Command())
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			userName := update.Message.From.UserName

			switch update.Message.Text {
			case "/start":
				if this.isAdmin(userName) {
					this.page = pages.NewAdminStartPage(this.db, this.adminQueryBuilder)
				} else {
					this.page = pages.NewStartPage(this.db, this.adminQueryBuilder)
				}
				msg.Text = this.page.RespText(update)
				msg.ReplyMarkup = this.page.(models.IPageWithKeyboard).Keyboard()
			default:
				if this.page.AllowedOnlyCommands() {
					msg.Text = "Select one option from the list."
					msg.ReplyMarkup = this.page.(models.IPageWithKeyboard).Keyboard()
				} else {
					pageWithAction, ok := this.page.(models.IPageWithAction)
					if ok {
						pageWithAction.Action(update)
					}

					nextPage := this.page.NextPage(update, this.isAdmin(userName))
					msg.Text = nextPage.RespText(update)
					nextPageWithKeyboard, ok := nextPage.(models.IPageWithKeyboard)
					if ok {
						msg.ReplyMarkup = nextPageWithKeyboard.Keyboard()
					}

					this.page = nextPage
				}
			}

			this.bot.Send(msg)
			// this.bot.SendMediaGroup()
			// this.bot.GetFile(tgbotapi.FileConfig{FileID: update.Message.Animation.FileID})
		} else if update.CallbackQuery != nil && update.CallbackQuery.Data != this.lastCommand {
			fmt.Printf("[%s] Data - %s\n", update.CallbackQuery.From.UserName, update.CallbackData())
			this.lastCommand = update.CallbackQuery.Data
			userName := update.CallbackQuery.From.UserName

			pageWithAction, ok := this.page.(models.IPageWithAction)
			if ok {
				pageWithAction.Action(update)
			}

			nextPage := this.page.NextPage(update, this.isAdmin(userName))
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, nextPage.RespText(update))

			nextPageWithKB, ok := nextPage.(models.IPageWithKeyboard)
			if ok {
				msg.ReplyMarkup = nextPageWithKB.Keyboard()
			}

			this.bot.Send(msg)
			this.page = nextPage
		}
	}
}

func (this *BuisdevBot) isAdmin(userName string) bool {
	for _, adminName := range this.admins {
		if strings.ToLower(userName) == strings.ToLower(adminName) {
			return true
		}
	}

	return false
}
