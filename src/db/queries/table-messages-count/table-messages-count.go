package t_msgs_count

import (
	"database/sql"

	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type T_MessagesCount struct {
	conn *sql.DB
}

func NewTableMessagesCount(conn *sql.DB) models.ITableMessagesCount {
	return T_MessagesCount{conn: conn}
}

func (this T_MessagesCount) CreateTable() error {
	_, err := this.conn.Exec(
		`CREATE TABLE IF NOT EXISTS messages_count (
            id INTEGER NOT NULL PRIMARY KEY,
            created_at TMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`,
	)

	return err
}

func (this T_MessagesCount) AddMessage(msg models.UserMsgFromClient) error {
	_, err := this.conn.Exec("INSERT INTO messages_count (created_at) VALUES ($1)", msg.CreatedAt)

	return err
}

func (this T_MessagesCount) CheckMessagesCount(fromTimestamp string) (int, error) {
	var count int
	err := this.conn.QueryRow(`
		SELECT COUNT(id) FROM messages_count
		WHERE created_at > $1;
	`, fromTimestamp).Scan(&count)

	return count, err
}

var _ models.ITableMessagesCount = &T_MessagesCount{}
