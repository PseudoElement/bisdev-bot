package main

import (
	"fmt"

	"github.com/joho/godotenv"
	bd_bot "github.com/pseudoelement/rubic-buisdev-tg-bot/src/bot"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	db := db.NewSqliteDB()

	bot := bd_bot.NewBuisdevBot(db)
	bot.Listen()
}
