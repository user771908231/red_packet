package model

import (
	"time"
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
	"casino_common/utils/timeUtils"
)

type MjRecordBean struct {
	UserId    uint32
	NickName  string
	WinAmount int64
}

func (b MjRecordBean) TransBeanUserRecord() *mjproto.BeanUserRecord {
	result := newProto.NewBeanUserRecord()
	*result.NickName = b.NickName
	*result.UserId = b.UserId
	*result.WinAmount = b.WinAmount
	return result
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

func (t T_mj_desk_round) TransRecord() *mjproto.BeanGameRecord {
	result := newProto.NewBeanGameRecord()
	*result.BeginTime = timeUtils.Format(t.BeginTime)
	*result.DeskId = t.DeskId
	for _, bean := range t.Records {
		b := bean.TransBeanUserRecord()
		result.Users = append(result.Users, b)
	}
	return result
}
