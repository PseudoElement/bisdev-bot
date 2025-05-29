package models

type DB_UserMessage struct {
	Id        int
	CreatedAt string
	UserName  string
	Initials  string
	Text      string
	New       bool
	Blob      []byte
	BlobType  string
}

type DB_UserNames struct {
	NotRead     []string
	AlreadyRead []string
}
