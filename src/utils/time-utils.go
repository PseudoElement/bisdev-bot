package utils

import "time"

// past arg means to create timestamp + mins or timestamp - mins
func GetSqlTimestampByMinutesUTC(mins int, past bool) string {
	var t time.Time
	if past {
		t = time.Now().UTC().Add(time.Duration(-mins) * time.Minute)
	} else {
		t = time.Now().UTC().Add(time.Duration(mins) * time.Minute)
	}
	formatted := t.Format("2006-01-02 15:04:05")

	return formatted
}

// past arg means to create timestamp + mins or timestamp - mins
func GetSqlTimestampByMinutes(mins int, past bool) string {
	var t time.Time
	if past {
		t = time.Now().Add(time.Duration(-mins) * time.Minute)
	} else {
		t = time.Now().Add(time.Duration(mins) * time.Minute)
	}
	formatted := t.Format("2006-01-02 15:04:05")

	return formatted
}
