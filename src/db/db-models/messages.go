package db_models

type DB_ClientMessage struct {
	Id        int
	CreatedAt string
	UserName  string
	Text      string
	New       bool
}
