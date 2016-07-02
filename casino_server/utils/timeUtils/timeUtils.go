package timeUtils

import "time"

const TIME_LAYOUT = "2006-01-02 15:04:05"

func Format(t time.Time) string{
	return t.Format(TIME_LAYOUT)

}
