package time

import (
	"time"
)

const (
	SecondsPerDay = 24 * 60 * 60

	UTCToBeijing = 8 * 60 * 60
	BeijingToUTC = -UTCToBeijing
)

var BeijingLocation = time.FixedZone("Asia/Shanghai", UTCToBeijing)

// UTC unixtime 转换为 北京时间 unixtime.
func UTCUnixToBeijingUnix(n int64) int64 {
	return n + UTCToBeijing
}

// 北京时间 unixtime 转换为 UTC unixtime.
func BeijingUnixToUTCUnix(n int64) int64 {
	return n + BeijingToUTC
}

// time.Time 转换为 北京时间 unixtime.
func TimeToBeijingUnix(t time.Time) int64 {
	return UTCUnixToBeijingUnix(t.Unix())
}

// 北京时间 unixtime 转换为 time.Time, BeijingLocation.
func BeijingUnixToTime(n int64) time.Time {
	return time.Unix(BeijingUnixToUTCUnix(n), 0).In(BeijingLocation)
}

//======================================================================================================================

// 北京时间 unixtime 转换为 北京时间距 1970-01-01 的天数.
func BeijingUnixToBeijingUnixDay(n int64) int64 {
	return n / SecondsPerDay
}

// 北京时间距 1970-01-01 的天数 转换为 北京时间 unixtime.
func BeijingUnixDayToBeijingUnix(n int64) int64 {
	return n * SecondsPerDay
}

// UTC unixtime 转换为 北京时间距 1970-01-01 的天数.
func UTCUnixToBeijingUnixDay(n int64) int64 {
	return BeijingUnixToBeijingUnixDay(UTCUnixToBeijingUnix(n))
}

// 北京时间距 1970-01-01 的天数 转换为 UTC unixtime.
func BeijingUnixDayToUTCUnix(n int64) int64 {
	return BeijingUnixToUTCUnix(BeijingUnixDayToBeijingUnix(n))
}

// time.Time 转换为 北京时间距 1970-01-01 的天数.
func TimeToBeijingUnixDay(t time.Time) int64 {
	return BeijingUnixToBeijingUnixDay(TimeToBeijingUnix(t))
}

// 北京时间距 1970-01-01 的天数 转换为 time.Time, BeijingLocation.
func BeijingUnixDayToTime(n int64) time.Time {
	return BeijingUnixToTime(BeijingUnixDayToBeijingUnix(n))
}

//======================================================================================================================

// UTC unixtime 转换为 UTC时间距 1970-01-01 的天数.
func UTCUnixToUTCUnixDay(n int64) int64 {
	return n / SecondsPerDay
}

// UTC时间距 1970-01-01 的天数 转换为 UTC unixtime.
func UTCUnixDayToUTCUnix(n int64) int64 {
	return n * SecondsPerDay
}

// time.Time 转换为 UTC时间距 1970-01-01 的天数.
func TimeToUTCUnixDay(t time.Time) int64 {
	return UTCUnixToUTCUnixDay(t.Unix())
}

// UTC时间距 1970-01-01 的天数 转换为 time.Time.
func UTCUnixDayToTime(n int64) time.Time {
	return time.Unix(UTCUnixDayToUTCUnix(n), 0)
}
