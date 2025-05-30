package models

type DB_UserMessage struct {
	Id        int
	CreatedAt string
	UserId    int64
	// username is update.UserName or update.UserID(if username doesn't exist)
	UserName string
	Initials string
	Text     string
	New      bool
	Blob     []byte
	BlobType string
}

type DB_UserNames struct {
	NotRead     []string
	AlreadyRead []string
}

type Db_BlockedUser struct {
	UserId    int64
	UserName  string
	CreatedAt string
}
