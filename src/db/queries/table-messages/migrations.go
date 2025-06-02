package t_messages

import (
	"database/sql"
	"fmt"

	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
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
	_, err := this.conn.Exec(fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS messages (
            id INTEGER NOT NULL PRIMARY KEY,
			user_id INTEGER NOT NULL,
            user_name VARCHAR(50) NOT NULL,
			initials VARCHAR(50) NOT NULL,
            text TEXT NOT NULL,
            new BOOLEAN NOT NULL,
			file_type VARCHAR(50) CHECK(file_type IN ('%s', '%s', '%s', '%s', '%s', '%s', '')) DEFAULT '',
			file_id VARCHAR(150) DEFAULT '',
            created_at TMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`, consts.FILE_TYPES.Doc, consts.FILE_TYPES.Csv, consts.FILE_TYPES.Txt, consts.FILE_TYPES.Pdf, consts.FILE_TYPES.Jpeg, consts.FILE_TYPES.Png),
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
