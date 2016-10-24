package model

import (
	"time"
)

type MjRecordBean struct {
	UserId    uint32
	NickName  string
	WinAmount int64
}


//一把结束,战绩可以通过这个表来查询
type T_mj_desk_round struct {
	DeskId     int32
	GameNumber int32
	UserIds    string
	BeginTime  time.Time
	EndTime    time.Time
	Records    []MjRecordBean
}