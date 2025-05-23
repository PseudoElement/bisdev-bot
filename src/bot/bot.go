package bd_bot

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages"
)

type BuisdevBot struct {
	bot         *tgbotapi.BotAPI
	isProd      bool
	page        models.IPage
	lastCommand string
	admins      []string
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

	b := &BuisdevBot{
		bot:         bot,
		isProd:      isProd,
		page:        pages.NewStartPage(),
		admins:      admins,
		lastCommand: "",
	}

	return b
}

func (this *BuisdevBot) ListenWithWebhook() {
	// pwd, _ := os.Getwd()
	// certFile := pwd + "/cert.pem"
	// keyFile := pwd + "/key.pem"

	// wh, err := tgbotapi.NewWebhookWithCert("https://amojo.amocrm.ru/~external/hooks/telegram?t="+this.bot.Token, tgbotapi.FilePath(certFile))
	// wh, err := tgbotapi.NewWebhook("https://amojo.amocrm.ru/~external/hooks/telegram?t=" + this.bot.Token)
	// wh, err := tgbotapi.NewWebhook("https://amojo.amocrm.ru/" + this.bot.Token)
	// _, err = this.bot.Request(wh)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	info, err := this.bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	// go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)

	log.Printf("Webhook ==> %+v \n", info)
	log.Printf("IsSet ==> %v \n", info.IsSet())
	log.Printf("Token ==> %+v \n", this.bot.Token)

	updatesChan := this.bot.ListenForWebhook(info.URL)
	// updatesChan := this.bot.ListenForWebhook("https://amojo.amocrm.ru/~external/hooks/telegram?t=" + this.bot.Token)

	for update := range updatesChan {
		fmt.Println("Updates")

		if update.Message != nil {
			fmt.Printf("[%s] webhook Message - %s\n", update.Message.From.UserName, update.Message.Text)
		}
		if update.CallbackQuery != nil {
			fmt.Printf("[%s] webhook CallbackQuery - %s\n", update.CallbackQuery.From.UserName, update.CallbackData())
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		this.bot.Send(msg)
	}
}

func (this *BuisdevBot) Listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range this.bot.GetUpdatesChan(u) {
		if update.Message != nil {
			// fmt.Printf("[%s] text - %s, command - %s.\n", update.Message.From.UserName, update.Message.Text, update.Message.Command())
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			switch update.Message.Text {
			case "/start":
				if this.isAdmin(update.Message.From.UserName) {
					this.page = pages.NewAdminStartPage()
				} else {
					this.page = pages.NewStartPage()
				}
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
			// fmt.Printf("[%s] Data - %s\n", update.CallbackQuery.From.UserName, update.CallbackQuery.Data)
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

func (this *BuisdevBot) isAdmin(userName string) bool {
	for _, adminName := range this.admins {
		if strings.ToLower(userName) == strings.ToLower(adminName) {
			return true
		}
	}

	return false
}
