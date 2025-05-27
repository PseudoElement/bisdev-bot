package quieres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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
			blob BLOB,
            created_at TMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`,
	)

	return err
}

func (this T_Messages) AddMessage(msg models.JsonMsgFromClient) error {
	log.Printf("[T_Messages_AddMessages] msg ==> %+v", msg)

	var err error
	if msg.ImageBlob != nil && len(msg.ImageBlob) > 0 {
		_, err = this.conn.Exec(
			"INSERT INTO messages (user_name, text, new, blob) VALUES ($1, $2, $3, $4)",
			msg.UserName, msg.Text, true, msg.ImageBlob)
	} else {
		_, err = this.conn.Exec(
			"INSERT INTO messages (user_name, text, new) VALUES ($1, $2, $3)",
			msg.UserName, msg.Text, true)
	}

	return err
}

func (this T_Messages) GetMessages(req models.MessagesReq) ([]models.DB_UserMessage, error) {
	messages := make([]models.DB_UserMessage, 0, req.Count)

	query := "SELECT id, user_name, text, new, created_at FROM messages "
	if req.NewOnly {
		query += "WHERE new = true "
	}
	query += "ORDER BY created_at DESC LIMIT $1;"
	log.Println("[GetMessages] query ==> ", query)

	rows, err := this.conn.Query(query, req.Count)
	if err != nil {
		return messages, err
	}
	defer rows.Close()

	for rows.Next() {
		msg := models.DB_UserMessage{}
		err := rows.Scan(&msg.Id, &msg.UserName, &msg.Text, &msg.New, &msg.CreatedAt)
		if err != nil {
			return messages, err
		}

		messages = append(messages, msg)
	}

	if len(messages) > 0 {
		go this.markMessagesAsRead(messages)
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

// LIMIT 10
func (this T_Messages) GetMessagesByUserName(userName string) ([]models.DB_UserMessage, error) {
	messages := make([]models.DB_UserMessage, 0, 5)

	rows, err := this.conn.Query(`
		SELECT * FROM messages WHERE user_name = $1 
		ORDER BY created_at DESC LIMIT 10;`,
		userName,
	)
	if err != nil {
		return messages, err
	}
	defer rows.Close()

	for rows.Next() {
		msg := models.DB_UserMessage{}
		err := rows.Scan(&msg.Id, &msg.UserName, &msg.Text, &msg.New, &msg.ImgBlob, &msg.CreatedAt)
		if err != nil {
			return messages, err
		}

		messages = append(messages, msg)
	}

	if len(messages) > 0 {
		go this.markMessagesAsRead(messages)
	}

	return messages, nil
}

func (this T_Messages) GetUserNames() (models.DB_UserNames, error) {
	userNames := models.DB_UserNames{
		NotRead:     make([]string, 10),
		AlreadyRead: make([]string, 10),
	}

	query := "SELECT * FROM messages;"
	rows, err := this.conn.Query(query)
	if err != nil {
		return userNames, err
	}
	defer rows.Close()

	for rows.Next() {
		msg := models.DB_UserMessage{}
		err := rows.Scan(&msg.Id, &msg.UserName, &msg.Text, &msg.New, &msg.ImgBlob, &msg.CreatedAt)
		if err != nil {
			return userNames, err
		}

		if msg.New {
			userNames.NotRead = append(userNames.NotRead, msg.UserName)
		} else {
			userNames.AlreadyRead = append(userNames.AlreadyRead, msg.UserName)
		}
	}

	userNames.AlreadyRead = this.uniqueUserNames(userNames.AlreadyRead)
	userNames.NotRead = this.uniqueUserNames(userNames.NotRead)

	return userNames, nil
}

func (this T_Messages) uniqueUserNames(allUserNames []string) []string {
	uniqueUserNames := make([]string, 0, len(allUserNames))
	m := make(map[string]int8)
	for _, name := range allUserNames {
		_, ok := m[name]
		if !ok {
			m[name] = 0
			uniqueUserNames = append(uniqueUserNames, name)
		}
	}

	return uniqueUserNames
}

func (this T_Messages) markMessagesAsRead(messages []models.DB_UserMessage) {
	placeholders := make([]string, 0, len(messages))
	newIds := make([]any, 0, len(messages))
	for _, msg := range messages {
		if msg.New {
			newIds = append(newIds, msg.Id)
			placeholders = append(placeholders, "?")
		}
	}

	query := fmt.Sprintf("UPDATE messages SET new = false WHERE id IN (%s)", strings.Join(placeholders, ","))
	log.Println("[T_Messages_markMessagesAsRead] query ==> ", query)

	_, err := this.conn.Exec(query, newIds...)
	if err != nil {
		log.Println("[T_Messages_markMessagesAsRead] err ==>", err)
	}
}

var _ models.ITableMessages = &T_Messages{}
