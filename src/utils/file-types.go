package utils

import "github.com/pseudoelement/rubic-buisdev-tg-bot/src/consts"

func IsImg(blobType string) bool {
	for _, t := range consts.IMAGES_FILE_TYPES {
		if t == blobType {
			return true
		}
	}
	return false
}

func IsDoc(blobType string) bool {
	for _, t := range consts.DOC_FILE_TYPES {
		if t == blobType {
			return true
		}
	}
	return false
}
