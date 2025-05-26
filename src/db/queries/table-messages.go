package quieres

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/db/db_models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type T_Messages struct {
	conn *sql.DB
}

func NewTableMessages(conn *sql.DB) T_Messages {
	return T_Messages{conn: conn}
}

func (this T_Messages) CreateTable() error {
	_, err := this.conn.Exec(
		`CREATE TABLE IF NOT EXISTS messages(
            id SERIAL NOT NULL PRIMARY KEY,
            user_name VARCHAR(50) NOT NULL,
            text TEXT NOT NULL,
            new BOOLEAN NOT NULL,
            created_at TMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`,
	)

	return err
}

func (this T_Messages) AddMessage(msg models.JsonClientMsg) error {
	_, err := this.conn.Exec(
		"INSERT INTO messages (user_name, text, new) VALUES ($1, $2, $3)",
		msg.UserName, msg.Text, true)

	return err
}

func (this T_Messages) GetMesages(req models.MessagesReq) ([]models.JsonClientMsg, error) {
	messages := make([]models.JsonClientMsg, 0, req.Count)
	ids := make([]int, 0, req.Count)

	// start tx
	tx, err := this.conn.Begin()
	if err != nil {
		return messages, err
	}
	defer tx.Rollback()

	query := "SELECT * FROM messages "
	if req.NewOnly {
		query += "WHERE new = true "
	}
	query += "LIMIT $1 ORDER BY created_at DESC;"

	rows, err := tx.Query(query, req.Count)
	if err != nil {
		return messages, err
	}
	defer rows.Close()

	for rows.Next() {
		dbMsg := db_models.DB_ClientMessage{}
		err := rows.Scan(&dbMsg.Id, &dbMsg.UserName, &dbMsg.Text, &dbMsg.New, &dbMsg.CreatedAt)
		if err != nil {
			return messages, err
		}

		ids = append(ids, dbMsg.Id)
		messages = append(messages, models.JsonClientMsg{UserName: dbMsg.UserName, Text: dbMsg.Text})
	}

	// Mark messages as read
	if len(ids) > 0 {
		_, err = tx.Exec("UPDATE messages SET new = false WHERE id = ANY($1)", pq.Array(ids))
		if err != nil {
			return messages, err
		}
	}

	// end tx
	err = tx.Commit()
	if err != nil {
		return messages, fmt.Errorf("[GetMesages] tx.Commit error: %v\n", err)
	}

	return messages, nil
}
