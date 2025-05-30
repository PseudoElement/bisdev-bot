package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	t_blocked_users "github.com/pseudoelement/rubic-buisdev-tg-bot/src/db/queries/table-blocked-users"
	t_messages "github.com/pseudoelement/rubic-buisdev-tg-bot/src/db/queries/table-messages"
	t_msgs_count "github.com/pseudoelement/rubic-buisdev-tg-bot/src/db/queries/table-messages-count"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type SqliteDB struct {
	conn   *sql.DB
	tables models.Tables
}

func NewSqliteDB() *SqliteDB {
	db := &SqliteDB{}
	conn, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic("[NewSqliteDB] sql.Open err:" + err.Error())
	}

	db.conn = conn

	db.tables = models.Tables{
		Messages:      t_messages.NewTableMessages(conn),
		MessagesCount: t_msgs_count.NewTableMessagesCount(conn),
		BlockedUsers:  t_blocked_users.NewTableBlockedUsers(conn),
	}

	err = db.tables.Messages.CreateTable()
	if err != nil {
		panic("[NewSqliteDB] Messages_CreateTable err ==>" + err.Error())
	}
	err = db.tables.MessagesCount.CreateTable()
	if err != nil {
		panic("[NewSqliteDB] MessagesCount_CreateTable err ==>" + err.Error())
	}
	err = db.tables.BlockedUsers.CreateTable()
	if err != nil {
		panic("[NewSqliteDB] BlockedUsers_CreateTable err ==>" + err.Error())
	}

	log.Println("Sqlite successfully connected.")

	return db
}

func (this SqliteDB) Conn() *sql.DB {
	return this.conn
}

func (this SqliteDB) Tables() models.Tables {
	return this.tables
}

var _ models.IDatabase = (*SqliteDB)(nil)
