package utils

import "time"

func StringToTime(str string)time.Time {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, str, loc)//使用模板在对应时区转化为time.time类型
	return  theTime
}
