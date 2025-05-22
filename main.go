package main

import (
	"fmt"

	"github.com/joho/godotenv"
	bd_bot "github.com/pseudoelement/rubic-buisdev-tg-bot/src/bot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	bot := bd_bot.NewBuisdevBot()
	bot.Listen()
}
