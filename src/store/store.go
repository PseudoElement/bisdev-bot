package store

import (
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/db"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type Store struct {
	db           *db.SqliteDB
	blockedUsers []models.Db_BlockedUser
	// key is UserID
	admins map[int64]models.Admin
}

func NewStore(db *db.SqliteDB) *Store {
	s := &Store{
		db:           db,
		blockedUsers: make([]models.Db_BlockedUser, 0, 10),
		admins:       make(map[int64]models.Admin, 10),
	}
	s.initAdmins()
	s.UpdateBlockedUsersList()

	return s
}

// sets all admins UserIDs
func (this *Store) initAdmins() {
	adminsString, ok := os.LookupEnv("ADMINS")
	if !ok {
		adminsString = ""
	}
	admins := strings.Split(adminsString, " ")

	for _, admin := range admins {
		userId, _ := strconv.Atoi(admin)
		this.admins[int64(userId)] = models.Admin{UserId: int64(userId)}
	}
}

func (this *Store) UpdateBlockedUsersList() []models.Db_BlockedUser {
	blockedUsers, err := this.db.Tables().BlockedUsers.GetBlockedUsers()
	if err != nil {
		log.Println("[Store_UpdateBlockedUsersList] GetBlockedUsers_err ==>", err)
	}

	this.blockedUsers = blockedUsers

	return this.blockedUsers
}

func (this *Store) GetBlockedUsers() []models.Db_BlockedUser {
	return this.blockedUsers
}

// checks by UserID
func (this *Store) IsBlockedUserById(userId int64) bool {
	for _, user := range this.blockedUsers {
		if userId == user.UserId {
			return true
		}
	}
	return false
}

// checks by UserID
func (this *Store) IsAdminById(userId int64) bool {
	for _, admin := range this.admins {
		if userId == admin.UserId {
			return true
		}
	}

	return false
}

// checks bu UserName. works only if admin typed any message to bot
func (this *Store) IsAdminByName(userName string) bool {
	for _, admin := range this.admins {
		if userName == admin.UserName {
			return true
		}
	}

	return false
}

func (this *Store) SetAdminData(update tgbotapi.Update) {
	var userId int64
	var userName string
	var chatId int64
	if update.Message != nil {
		userId = update.Message.From.ID
		userName = update.Message.From.UserName
		chatId = update.Message.Chat.ID
	}
	if update.CallbackQuery != nil {
		userId = update.CallbackQuery.From.ID
		userName = update.CallbackQuery.From.UserName
		chatId = update.CallbackQuery.Message.Chat.ID
	}

	this.admins[userId] = models.Admin{
		ChatId:             chatId,
		UserId:             chatId,
		UserName:           userName,
		IsListenToNotifier: true,
	}
}

func (this *Store) GetAdmins() map[int64]models.Admin {
	return this.admins
}
