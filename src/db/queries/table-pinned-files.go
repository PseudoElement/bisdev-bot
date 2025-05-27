package quieres

import (
	"database/sql"

	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type T_PinnedFiles struct {
	conn *sql.DB
}

func NewTablePinnedFiles(conn *sql.DB) models.ITablePinnedFiles {
	return T_PinnedFiles{conn: conn}
}

func (this T_PinnedFiles) CreateTable() error {
	_, err := this.conn.Exec(
		`CREATE TABLE IF NOT EXISTS pinned_files (
            message_id INTEGER NOT NULL,
			file_type VARCHAR(50) CHECK( file_type IN ($1, $2) ),
			blob BLOB NOT NULL,
			FOREIGN KEY(message_id) REFERENCES messages(id)
        );`, consts.FILE_TYPES.Image, consts.FILE_TYPES.Document,
	)

	return err
}

// @TODO add queries to send for each message every pinned file for it
