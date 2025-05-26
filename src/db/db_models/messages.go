package db_models

type DB_ClientMessage struct {
	Id        int
	CreatedAt int
	UserName  string
	Text      string
	New       bool
}
