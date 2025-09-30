package t_admins

import (
	"database/sql"
	"log"

	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/models"
)

type T_Admins struct {
	conn *sql.DB
}

func NewTableAdmins(conn *sql.DB) models.ITableAdmins {
	return T_Admins{conn: conn}
}

func (this T_Admins) CreateTable() error {
	_, err := this.conn.Exec(
		`CREATE TABLE IF NOT EXISTS admins (
            user_id INTEGER NOT NULL,
            chat_id INTEGER NOT NULL,
            user_name VARCHAR(50) NOT NULL
        );`,
	)

	return err
}

func (this T_Admins) SaveAdmin(adminInfo models.DB_Admin) error {
	res, err := this.conn.Exec(`INSERT INTO admins VALUES($1, $2, $3)`, adminInfo.UserId, adminInfo.ChatId, adminInfo.UserName)
	if err != nil {
		log.Println("[T_Admins_SaveAdmin] Exec_err ==>", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Println("[T_Admins_SaveAdmin] RowsAffected_err ==>", err)
	}

	log.Println("[T_Admins_SaveAdmin] count ==>", count)

	return err
}

func (this T_Admins) GetAdmins() ([]models.DB_Admin, error) {
	admins := make([]models.DB_Admin, 0, 10)

	rows, err := this.conn.Query("SELECT * FROM admins;")
	if err != nil {
		return admins, err
	}
	defer rows.Close()

	for rows.Next() {
		admin := models.DB_Admin{}
		err := rows.Scan(&admin.UserId, &admin.ChatId, &admin.UserName)
		if err != nil {
			return admins, err
		}

		admins = append(admins, admin)
	}

	return admins, nil
}
