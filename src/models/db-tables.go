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
	Admins        ITableAdmins
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

	GetUserId(userName string) int64
}

type ITableMessagesCount interface {
	CreateTable() error

	AddMessage(msg UserMsgFromClient) error

	CheckMessagesCount(fromTimestamp string) (int, error)
}

type ITableBlockedUsers interface {
	CreateTable() error

	BlockUser(userName string) error

	UnblockUser(userName string) error

	GetBlockedUsers() ([]Db_BlockedUser, error)
}

type ITableAdmins interface {
	CreateTable() error

	SaveAdmin(adminInfo DB_Admin) error

	GetAdmins() ([]DB_Admin, error)
}

// @IMPROVEMENT add opportunity to store many messages
type ITablePinnedFiles interface{}
