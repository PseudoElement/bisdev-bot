package utils

const MAX_RUNE_LEN = 4_000

func SplitLongTextForTg(text string) []string {
	runes := []rune(text)
	length := len(runes)

	if length <= MAX_RUNE_LEN {
		return []string{text}
	}

	var chunks []string
	for i := 0; i < length; i += MAX_RUNE_LEN {
		end := i + MAX_RUNE_LEN
		if end > length {
			end = length
		}
		chunks = append(chunks, string(runes[i:end]))
	}

	return chunks
}
