package bd_bot

import (
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages"
)

type BuisdevBot struct {
	bot         *tgbotapi.BotAPI
	isProd      bool
	page        models.IPage
	lastCommand string
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

	bot.Debug = !isProd

	b := &BuisdevBot{
		bot:    bot,
		isProd: isProd,
		page:   pages.NewStartPage(),
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return b
}

func (this *BuisdevBot) Listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range this.bot.GetUpdatesChan(u) {
		if update.Message != nil {
			fmt.Printf("[%s] text - %s, command - %s.\n", update.Message.From.UserName, update.Message.Text, update.Message.Command())

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			switch update.Message.Text {
			case "/start":
				this.page = pages.NewStartPage()
				msg.Text = this.page.RespText(update)
				msg.ReplyMarkup = this.page.(models.IPageWithKeyboard).Keyboard()
			default:
				if this.page.AllowedOnlyCommands() {
					msg.Text = "Select one option from the list."
					msg.ReplyMarkup = this.page.(models.IPageWithKeyboard).Keyboard()
				} else {
					nextPage := this.page.NextPage(update)
					msg.Text = nextPage.RespText(update)
					msg.ReplyMarkup = nextPage.(models.IPageWithKeyboard).Keyboard()
					this.page = nextPage
				}
			}

			this.bot.Send(msg)
		} else if update.CallbackQuery != nil && update.CallbackQuery.Data != this.lastCommand {
			fmt.Printf("[%s] Data - %s\n", update.CallbackQuery.From.UserName, update.CallbackQuery.Data)
			this.lastCommand = update.CallbackQuery.Data

			nextPage := this.page.NextPage(update)
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
