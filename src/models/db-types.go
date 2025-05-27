package models

import "database/sql"

type IDatabase interface {
	Conn() *sql.DB

	Tables() Tables
}

type Tables struct {
	Messages    ITableMessages
	PinnedFiles ITablePinnedFiles
}

type ITableMessages interface {
	CreateTable() error

	AddMessage(msg JsonMsgFromClient) error

	GetMessages(req MessagesReq) ([]DB_UserMessage, error)

	GetMessagesByUserName(userName string) ([]DB_UserMessage, error)

	DeleteMessages(count int) error

	GetUserNames() (DB_UserNames, error)
}

type ITablePinnedFiles interface{}
