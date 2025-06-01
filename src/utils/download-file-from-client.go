package utils

import (
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetBytesByTgFileID(bot *tgbotapi.BotAPI, fileID string) ([]byte, error) {
	// Get file info
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return nil, err
	}

	// Get full file URL
	url := file.Link(bot.Token)

	// Fetch file
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read bytes
	return io.ReadAll(resp.Body)
}
