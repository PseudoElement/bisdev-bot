package quieres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	db_models "github.com/pseudoelement/rubic-buisdev-tg-bot/src/db/db-models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type T_Messages struct {
	conn *sql.DB
}

func NewTableMessages(conn *sql.DB) models.ITableMessages {
	return T_Messages{conn: conn}
}

func (this T_Messages) CreateTable() error {
	_, err := this.conn.Exec(
		`CREATE TABLE IF NOT EXISTS messages (
            id INTEGER NOT NULL PRIMARY KEY,
            user_name VARCHAR(50) NOT NULL,
            text TEXT NOT NULL,
            new BOOLEAN NOT NULL,
            created_at TMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`,
	)

	return err
}

func (this T_Messages) AddMessage(msg models.JsonClientMsg) error {
	log.Printf("[T_Messages_AddMessages] msg ==> %+v", msg)
	_, err := this.conn.Exec(
		"INSERT INTO messages (user_name, text, new) VALUES ($1, $2, $3)",
		msg.UserName, msg.Text, true)

	return err
}

func (this T_Messages) GetMesages(req models.MessagesReq) ([]models.JsonClientMsg, error) {
	messages := make([]models.JsonClientMsg, 0, req.Count)
	dbMessages := make([]db_models.DB_ClientMessage, 0, req.Count)

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
	query += "ORDER BY created_at DESC LIMIT $1;"
	log.Println("[GetMessages] query ==> ", query)

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

		dbMessages = append(dbMessages, dbMsg)
		messages = append(messages, models.JsonClientMsg{UserName: dbMsg.UserName, Text: dbMsg.Text})
	}

	// Mark messages as read
	go func() {
		if len(dbMessages) > 0 {
			placeholders := make([]string, 0, len(dbMessages))
			newIds := make([]any, 0, len(dbMessages))
			for _, dbMsg := range dbMessages {
				if dbMsg.New {
					newIds = append(newIds, dbMsg.Id)
					placeholders = append(placeholders, "?")
				}
			}

			query := fmt.Sprintf("UPDATE messages SET new = false WHERE id IN (%s)", strings.Join(placeholders, ","))
			log.Println("[T_Messages_GetMessages] update query ==> ", query)

			_, err = this.conn.Exec(query, newIds...)
			if err != nil {
				log.Println("[T_Messages_GetMessages] update to false ==>", err)
			}
		}
	}()

	// end tx
	err = tx.Commit()
	if err != nil {
		return messages, fmt.Errorf("[GetMesages] tx.Commit error: %v\n", err)
	}

	return messages, nil
}

func (this T_Messages) DeleteMessages(count int) error {
	_, err := this.conn.Exec(`
		DELETE FROM messages WHERE id IN (
			SELECT id FROM messages ORDER BY created_at ASC LIMIT $1
		);`, count)
	return err
}
