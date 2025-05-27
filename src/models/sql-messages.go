package models

type DB_UserMessage struct {
	Id        int
	CreatedAt string
	UserName  string
	Text      string
	New       bool
	ImgBlob   []byte `json:"img_blob,omitempty"`
}

type DB_UserNames struct {
	NotRead     []string
	AlreadyRead []string
}
