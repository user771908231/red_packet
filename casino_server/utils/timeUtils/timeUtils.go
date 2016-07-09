package timeUtils

import "time"

const TIME_LAYOUT = "2006-01-02 15:04:05"
const TIME_LAYOUT_YYYY_MM_DD = "2006-01-02"

func Format(t time.Time) string {
	return t.Format(TIME_LAYOUT)
}

func FormatYYYYMMDD(t time.Time) string {
	return t.Format(TIME_LAYOUT_YYYY_MM_DD)
}
