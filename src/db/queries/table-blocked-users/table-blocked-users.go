package t_blocked_users

import (
	"database/sql"
	"log"

	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type T_BlockedUsers struct {
	conn *sql.DB
}

func NewTableBlockedUsers(conn *sql.DB) models.ITableBlockedUsers {
	return T_BlockedUsers{conn: conn}
}

func (this T_BlockedUsers) CreateTable() error {
	_, err := this.conn.Exec(
		`CREATE TABLE IF NOT EXISTS blocked_users (
            user_id INTEGER NOT NULL,
			user_name VARCHAR(50) NOT NULL,
			created_at TMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`,
	)

	return err
}

func (this T_BlockedUsers) BlockUser(userName string) error {
	res, err := this.conn.Exec(`
		INSERT INTO blocked_users (user_id, user_name)
		SELECT (user_id, user_name) FROM messages
		WHERE user_name = $1;
	`, userName)
	if err != nil {
		log.Println("[T_BlockedUsers_BlockUser] Exec_err ==>", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Println("[T_BlockedUsers_BlockUser] RowsAffected_err ==>", err)
	}

	log.Println("[T_BlockedUsers_BlockUser] blocked_count ==>", count)

	return err
}

func (this T_BlockedUsers) GetBlockedUsers() ([]models.Db_BlockedUser, error) {
	blockedUsers := make([]models.Db_BlockedUser, 0, 5)

	rows, err := this.conn.Query("SELECT * FROM blocked_users;")
	if err != nil {
		return blockedUsers, err
	}
	defer rows.Close()

	for rows.Next() {
		user := models.Db_BlockedUser{}
		err := rows.Scan(&user.UserId, &user.UserName, &user.CreatedAt)
		if err != nil {
			return blockedUsers, err
		}

		blockedUsers = append(blockedUsers, user)
	}

	return blockedUsers, nil
}

var _ models.ITableBlockedUsers = &T_BlockedUsers{}
