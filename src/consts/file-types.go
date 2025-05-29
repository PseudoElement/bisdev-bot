package consts

type fileTypes struct {
	Doc  string
	Csv  string
	Txt  string
	Pdf  string
	Png  string
	Jpeg string
}

// change format
var FILE_TYPES = fileTypes{
	Doc:  "doc",
	Csv:  "csv",
	Txt:  "txt",
	Pdf:  "pdf",
	Png:  "png",
	Jpeg: "jpeg",
}

var SUPPORTED_MIME_TYPES = [6]string{
	"image/jpeg",
	"image/png",
	"application/pdf",
	"application/msword",
	"text/plain",
	"text/csv",
}

var IMAGES_FILE_TYPES = [2]string{
	FILE_TYPES.Jpeg,
	FILE_TYPES.Png,
}

var DOC_FILE_TYPES = [4]string{
	FILE_TYPES.Doc,
	FILE_TYPES.Csv,
	FILE_TYPES.Txt,
	FILE_TYPES.Pdf,
}
