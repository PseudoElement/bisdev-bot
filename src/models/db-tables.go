package models

import "database/sql"

type IDatabase interface {
	Conn() *sql.DB

	Tables() Tables
}

type Tables struct {
	Messages      ITableMessages
	MessagesCount ITableMessagesCount
	BlockedUsers  ITableBlockedUsers
	// PinnedFiles   ITablePinnedFiles
}

type ITableMessages interface {
	CreateTable() error

	AlterTable(columnName string, defaultValue string) error

	AddMessage(msg UserMsgFromClient) error

	GetMessages(req MessagesReq) ([]DB_UserMessage, error)

	GetMessagesByUserName(userName string) ([]DB_UserMessage, error)

	DeleteMessages(count int) error

	DeleteMessagesByUserName(userName string) error

	GetUserNames() (DB_UserNames, error)
}

type ITableMessagesCount interface {
	CreateTable() error

	AddMessage(msg UserMsgFromClient) error

	CheckMessagesCount(fromTimestamp string) (int, error)
}

type ITableBlockedUsers interface {
	CreateTable() error

	BlockUser(userName string) error

	GetBlockedUsers() ([]Db_BlockedUser, error)
}

// @IMPROVEMENT add opportunity to store many messages
type ITablePinnedFiles interface{}
