package db

import (
	"database/sql"

	quieres "github.com/pseudoelement/rubic-buisdev-tg-bot/src/db/queries"
)

type Tables struct {
	Messages quieres.T_Messages
}

type SqliteDB struct {
	Conn   *sql.DB
	tables Tables
}

func NewSqliteDB() *SqliteDB {
	db := &SqliteDB{}
	conn, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic("[error] sql.Open:" + err.Error())
	}

	db.Conn = conn

	db.tables = Tables{
		Messages: quieres.NewTableMessages(conn),
	}

	return db
}

func (this SqliteDB) Tables() Tables {
	return this.tables
}
