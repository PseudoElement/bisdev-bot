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

	AllowedOnlyMessages() bool
}

type IPageWithKeyboard interface {
	IPage

	Keyboard() tgbotapi.InlineKeyboardMarkup
}

type IPageWithActionOnDestroy interface {
	IPage

	ActionOnDestroy(update tgbotapi.Update)
}

type IPageWithActionOnInit interface {
	IPage

	ActionOnInit(update tgbotapi.Update)
}

type IPageWithPhotos interface {
	IPage

	PhotosResp(update tgbotapi.Update) tgbotapi.MediaGroupConfig

	HasPhotos() bool
}

type IPageWithFiles interface {
	IPage

	FilesResp(update tgbotapi.Update) tgbotapi.MediaGroupConfig

	HasFiles() bool
}
