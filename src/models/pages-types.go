package models

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type IPage interface {
	RespText(update tgbotapi.Update) string
	/* handles selected command and send next answer with buttons(or not) */
	NextPage(update tgbotapi.Update, isAdmin bool) IPage

	Name() string

	AllowedOnlyCommands() bool
}

type IPageWithKeyboard interface {
	IPage

	Keyboard() tgbotapi.InlineKeyboardMarkup
}

type IPageWithAction interface {
	IPage

	Action(update tgbotapi.Update)
}
