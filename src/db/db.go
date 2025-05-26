package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	quieres "github.com/pseudoelement/rubic-buisdev-tg-bot/src/db/queries"
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
		Messages: quieres.NewTableMessages(conn),
	}

	err = db.tables.Messages.CreateTable()
	if err != nil {
		panic("[NewSqliteDB] CreateTable err:" + err.Error())
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
