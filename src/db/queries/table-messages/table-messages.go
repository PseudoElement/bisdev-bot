package t_messages

import (
	"fmt"
	"log"
	"strings"

	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/utils"
)

func (this T_Messages) AddMessage(msg models.JsonMsgFromClient) error {
	log.Printf("[T_Messages_AddMessages] msg ==> %+v", models.MsgFromClientForLog{msg.UserName, msg.Initials, msg.Text, len(msg.Blob), msg.BlobType})

	var err error
	if msg.Blob != nil && len(msg.Blob) > 0 {
		_, err = this.conn.Exec(
			`INSERT INTO messages (user_name, initials, text, new, blob_type, blob, created_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7);`,
			msg.UserName, msg.Initials, msg.Text, true, msg.BlobType, msg.Blob, msg.CreatedAt)
	} else {
		_, err = this.conn.Exec(
			`INSERT INTO messages (user_name, initials, text, new, created_at) 
			VALUES ($1, $2, $3, $4, $5);`,
			msg.UserName, msg.Initials, msg.Text, true, msg.CreatedAt)
	}

	return err
}

func (this T_Messages) GetMessages(req models.MessagesReq) ([]models.DB_UserMessage, error) {
	messages := make([]models.DB_UserMessage, 0, req.Count)
	hasNewMsgs := false

	query := "SELECT id, user_name, initials, text, new, created_at FROM messages "
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
		err := rows.Scan(&msg.Id, &msg.UserName, &msg.Initials, &msg.Text, &msg.New, &msg.CreatedAt)
		if err != nil {
			return messages, err
		}
		if msg.New {
			hasNewMsgs = true
		}

		messages = append(messages, msg)
	}

	if len(messages) > 0 && hasNewMsgs {
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

func (this T_Messages) DeleteMessagesByUserName(userName string) error {
	_, err := this.conn.Exec("DELETE FROM messages WHERE user_name = $1;", userName)
	return err
}

// LIMIT 10
func (this T_Messages) GetMessagesByUserName(userName string) ([]models.DB_UserMessage, error) {
	messages := make([]models.DB_UserMessage, 0, 5)
	hasNewMsgs := false

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
		err := rows.Scan(&msg.Id, &msg.UserName, &msg.Initials, &msg.Text, &msg.New, &msg.BlobType, &msg.Blob, &msg.CreatedAt)
		if err != nil {
			return messages, err
		}
		if msg.New {
			hasNewMsgs = true
		}

		messages = append(messages, msg)
	}

	if len(messages) > 0 && hasNewMsgs {
		go this.markMessagesAsRead(messages)
	}

	return messages, nil
}

func (this T_Messages) GetUserNames() (models.DB_UserNames, error) {
	userNames := models.DB_UserNames{
		NotRead:     make([]string, 0, 10),
		AlreadyRead: make([]string, 0, 10),
	}

	query := "SELECT user_name, new FROM messages;"
	rows, err := this.conn.Query(query)
	if err != nil {
		return userNames, err
	}
	defer rows.Close()

	for rows.Next() {
		msg := models.DB_UserMessage{}
		err := rows.Scan(&msg.UserName, &msg.New)
		if err != nil {
			return userNames, err
		}

		if msg.New {
			userNames.NotRead = append(userNames.NotRead, msg.UserName)
		} else {
			userNames.AlreadyRead = append(userNames.AlreadyRead, msg.UserName)
		}
	}

	userNames.AlreadyRead = utils.FilterUnique(userNames.AlreadyRead)
	userNames.NotRead = utils.FilterUnique(userNames.NotRead)

	return userNames, nil
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
