package t_messages

import (
	"database/sql"
	"fmt"

	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type T_Messages struct {
	conn *sql.DB
}

func NewTableMessages(conn *sql.DB) models.ITableMessages {
	return T_Messages{conn: conn}
}

// UP
func (this T_Messages) CreateTable() error {
	_, err := this.conn.Exec(
		`CREATE TABLE IF NOT EXISTS messages (
            id INTEGER NOT NULL PRIMARY KEY,
            user_name VARCHAR(50) NOT NULL,
			initials VARCHAR(50) NOT NULL,
            text TEXT NOT NULL,
            new BOOLEAN NOT NULL,
			blob BLOB,
            created_at TMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`,
	)

	return err
}

func (this T_Messages) AlterTable(columnName string, defaultValue string) error {
	query := fmt.Sprintf(
		`ALTER TABLE messages ADD COLUMN %s VARCHAR(50) NOT NULL DEFAULT '%s';`,
		columnName, defaultValue,
	)
	_, err := this.conn.Exec(query)

	return err
}
