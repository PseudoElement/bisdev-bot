package injector

import (
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/db"
	query_builder "github.com/pseudoelement/rubic-buisdev-tg-bot/src/query-builder"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/store"
)

type AppInjector struct {
	Db                *db.SqliteDB
	Store             *store.Store
	AdminQueryBuilder *query_builder.AdminQueryBuilder
}

func NewAppInjector() *AppInjector {
	db := db.NewSqliteDB()
	adminQueryBuilder := query_builder.NewAdminQueryBuilder()
	store := store.NewStore(db)

	i := &AppInjector{
		AdminQueryBuilder: adminQueryBuilder,
		Db:                db,
		Store:             store,
	}

	return i
}
