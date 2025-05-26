package models

import "database/sql"

type IDatabase interface {
	Conn() *sql.DB

	Tables() Tables
}

type Tables struct {
	Messages ITableMessages
}

type ITableMessages interface {
	CreateTable() error

	AddMessage(msg JsonClientMsg) error

	GetMesages(req MessagesReq) ([]JsonClientMsg, error)

	DeleteMessages(count int) error
}
