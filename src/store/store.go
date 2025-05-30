package store

import (
	"log"

	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/db"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type Store struct {
	db           *db.SqliteDB
	blockedUsers []models.Db_BlockedUser
	admins       []models.Admin
}

func NewStore(db *db.SqliteDB) *Store {
	s := &Store{
		db:           db,
		blockedUsers: make([]models.Db_BlockedUser, 0, 10),
	}
	s.LoadBlockedUsers()

	return s
}

func (this *Store) LoadBlockedUsers() []models.Db_BlockedUser {
	blockedUsers, err := this.db.Tables().BlockedUsers.GetBlockedUsers()
	if err != nil {
		log.Println("[Store_LoadBlockedUsers] GetBlockedUsers_err ==>", err)
	}

	this.blockedUsers = blockedUsers

	return this.blockedUsers
}

func (this *Store) GetBlockedUsers() []models.Db_BlockedUser {
	return this.blockedUsers
}
