package utils

import (
	"strings"

	"github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"
)

func MimeTypeToSqlBlobType(mimetype string) string {
	switch mimetype {
	case "image/jpeg":
		return consts.FILE_TYPES.Jpeg
	case "image/png":
		return consts.FILE_TYPES.Png
	case "application/pdf":
		return consts.FILE_TYPES.Pdf
	case "application/msword":
		return consts.FILE_TYPES.Doc
	case "text/plain":
		return consts.FILE_TYPES.Txt
	case "text/csv":
		return consts.FILE_TYPES.Csv
	default:
		return consts.FILE_TYPES.Txt
	}
}

func UserNameForSql(userName string) string {
	return strings.ToLower(userName)
}
