package bd_bot

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/injector"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/notifier"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/pages"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/utils"
)

type BuisdevBot struct {
	bot *tgbotapi.BotAPI
	// db                *db.SqliteDB
	isProd      bool
	lastCommand string
	// admins            []string
	// adminQueryBuilder *query_builder.AdminQueryBuilder
	// key is userId
	pages map[int64]models.IPage
	// store    *store.Store
	injector *injector.AppInjector
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

	injector := injector.NewAppInjector(bot)

	b := &BuisdevBot{
		bot:         bot,
		isProd:      isProd,
		pages:       make(map[int64]models.IPage, 10),
		lastCommand: "",
		injector:    injector,
	}

	return b
}

func (this *BuisdevBot) ListenUpdates() {
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

		if this.injector.Store.IsBlockedUserById(userId) && !this.injector.Store.IsAdminById(userId) {
			this.handleBlockedUserRequest(update)
			continue
		}

		_, ok := this.pages[userId]
		if !ok {
			if this.injector.Store.IsAdminById(userId) {
				this.pages[userId] = pages.NewAdminStartPage(this.injector)
			} else {
				this.pages[userId] = pages.NewStartPage(this.injector)
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

func (this *BuisdevBot) ListenNotifier() {
	for note := range this.injector.Notifier.Chan() {
		switch note.(type) {
		case notifier.NotificationBlockUser:
			v := note.(notifier.NotificationBlockUser)
			go this.sendBlockInfoToAdmins(v)
		case notifier.NotificationUnblockUser:
			v := note.(notifier.NotificationUnblockUser)
			go this.sendUnblockInfoToAdmins(v)
		case notifier.NotificationNewMessage:
			v := note.(notifier.NotificationNewMessage)
			go this.sendNewMessageToAdmins(v)
		default:
			log.Println("!!!NOTE: unknown note type %v", note)
		}
	}
}

func (this *BuisdevBot) handleMessageRequest(update tgbotapi.Update) {
	userId := update.Message.From.ID

	if this.injector.Store.IsAdminById(userId) {
		this.injector.Store.SetAdminData(update)
	}

	switch update.Message.Text {
	case "/start":
		this.lastCommand = ""
		if this.injector.Store.IsAdminById(userId) {
			this.pages[userId] = pages.NewAdminStartPage(this.injector)
		} else {
			this.pages[userId] = pages.NewStartPage(this.injector)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.Text = this.pages[userId].RespText(update)
		msg.ReplyMarkup = this.pages[userId].(models.IPageWithKeyboard).Keyboard()

		go this.bot.Send(msg)
	default:
		if update.Message.Text == "" && update.Message.Caption == "" {
			// ignored messages with media without text
			// ignored all media pinned except first
			return
		} else if this.pages[userId].AllowedOnlyCommands() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.Text = "Select one option from the list."
			msg.ReplyMarkup = this.pages[userId].(models.IPageWithKeyboard).Keyboard()

			go this.bot.Send(msg)
		} else {
			pageWithActionOnDestr, withOnDestroy := this.pages[userId].(models.IPageWithActionOnDestroy)
			if withOnDestroy {
				pageWithActionOnDestr.ActionOnDestroy(update)
			}

			nextPage := this.pages[userId].NextPage(update, this.injector.Store.IsAdminById(userId))
			nextPageWithActionOnInit, withOnInit := nextPage.(models.IPageWithActionOnInit)
			nextPageWithPhotos, withPhotosInResp := nextPage.(models.IPageWithPhotos)
			nextPageWithFiles, withFilesInResp := nextPage.(models.IPageWithFiles)

			if withOnInit {
				nextPageWithActionOnInit.ActionOnInit(update)
			}

			// @TODO add more than 10 files/photos per request
			go func() {
				this.sendTextResponse(update, nextPage)
				if withPhotosInResp && nextPageWithPhotos.HasPhotos() {
					photoGroups := nextPageWithPhotos.PhotosResp(update)
					for _, photoGroup := range photoGroups {
						this.bot.SendMediaGroup(photoGroup)
					}
				}
				if withFilesInResp && nextPageWithFiles.HasFiles() {
					photoGroups := nextPageWithFiles.FilesResp(update)
					for _, photoGroup := range photoGroups {
						this.bot.SendMediaGroup(photoGroup)
					}
				}
			}()

			this.pages[userId] = nextPage
		}
	}
}

func (this *BuisdevBot) handleCallbackRequest(update tgbotapi.Update) {
	userId := update.CallbackQuery.From.ID

	if this.injector.Store.IsAdminById(userId) {
		this.injector.Store.SetAdminData(update)
	}

	if this.pages[userId].AllowedOnlyMessages() && update.CallbackData() != consts.BACK_TO_START {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Commands are not allowed for response here.")

		go this.bot.Send(msg)
	} else {
		pageWithActionOnDestr, ok := this.pages[userId].(models.IPageWithActionOnDestroy)
		// skip onDestroy callback if want to get back to start
		if ok && update.CallbackData() != consts.BACK_TO_START {
			pageWithActionOnDestr.ActionOnDestroy(update)
		}

		nextPage := this.pages[userId].NextPage(update, this.injector.Store.IsAdminById(userId))
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

func (this *BuisdevBot) handleBlockedUserRequest(update tgbotapi.Update) {
	var chatId int64
	if update.Message != nil {
		chatId = update.Message.Chat.ID
	}
	if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.Message.Chat.ID
	}
	msg := tgbotapi.NewMessage(
		chatId,
		"You're blocked because of rules violation. Contact support-team for details https://t.me/RubicSupportBot.",
	)

	go this.bot.Send(msg)
}

// if response longer than 4000 characters - it splits it into chunks by 4000 chars per chunk
func (this *BuisdevBot) sendTextResponse(update tgbotapi.Update, nextPage models.IPage) {
	textChunks := utils.SplitLongTextForTg(nextPage.RespText(update))
	nextPageWithKeyboard, withKeyboard := nextPage.(models.IPageWithKeyboard)

	for idx, textChunk := range textChunks {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, textChunk)
		// send keyboard with last chunk if needed
		if idx == len(textChunks)-1 && withKeyboard {
			msg.ReplyMarkup = nextPageWithKeyboard.Keyboard()
		}

		this.bot.Send(msg)
	}
}

func (this *BuisdevBot) sendNewMessageToAdmins(note notifier.NotificationNewMessage) {
	for _, admin := range this.injector.Store.GetAdmins() {
		text := "✉️ New message from " + note.FromInitials + "(login @" + note.FromUserName + ")\n"
		if note.WithFiles {
			text += "Contains pinned files, to see files - load messages of this user via **👤 Show messages of specific user**.\n"
		}
		text += "Message:\n" + note.Text
		msg := tgbotapi.NewMessage(admin.ChatId, text)

		this.bot.Send(msg)
	}
}

func (this *BuisdevBot) sendBlockInfoToAdmins(note notifier.NotificationBlockUser) {
	for _, admin := range this.injector.Store.GetAdmins() {
		text := fmt.Sprintf("🚷 %s blocked user %s.", note.AdminUserName, note.BlockedUserName)
		msg := tgbotapi.NewMessage(admin.ChatId, text)

		this.bot.Send(msg)
	}
}

func (this *BuisdevBot) sendUnblockInfoToAdmins(note notifier.NotificationUnblockUser) {
	for _, admin := range this.injector.Store.GetAdmins() {
		text := fmt.Sprintf("💫 %s unblocked user %s.", note.AdminUserName, note.UnblockedUserName)
		msg := tgbotapi.NewMessage(admin.ChatId, text)

		this.bot.Send(msg)
	}
}
