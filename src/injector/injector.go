package injector

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/db"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/store"
)

type AppInjector struct {
	Db                *db.SqliteDB
	Store             *store.Store
	AdminQueryBuilder *query_builder.AdminQueryBuilder
	Bot               *tgbotapi.BotAPI
}

func NewAppInjector(bot *tgbotapi.BotAPI) *AppInjector {
	db := db.NewSqliteDB()
	adminQueryBuilder := query_builder.NewAdminQueryBuilder()
	store := store.NewStore(db)

	i := &AppInjector{
		AdminQueryBuilder: adminQueryBuilder,
		Db:                db,
		Store:             store,
		Bot:               bot,
	}

	return i
}
