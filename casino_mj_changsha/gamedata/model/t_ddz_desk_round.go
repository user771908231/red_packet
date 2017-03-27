package model

import (
	"time"
	"casino_mj_changsha/msg/funcsInit"
	"casino_common/utils/timeUtils"
	"casino_mj_changsha/msg/protogo"
)

type DdzRecordBean struct {
	UserId    uint32
	NickName  string
	WinAmount int64
}

func (b DdzRecordBean) TransBeanUserRecord() *csmjproto.BeanUserRecord {
	result := newProto.NewBeanUserRecord()
	*result.NickName = b.NickName
	*result.UserId = b.UserId
	*result.WinAmount = b.WinAmount
	return result
}



//一把结束,战绩可以通过这个表来查询
type T_ddz_desk_round struct {
	DeskId     int32
	GameNumber int32
	UserIds    string
	BeginTime  time.Time
	EndTime    time.Time
	Records    []DdzRecordBean
}

func (t T_ddz_desk_round) TransRecord() *csmjproto.BeanGameRecord {
	result := newProto.NewBeanGameRecord()
	*result.BeginTime = timeUtils.Format(t.BeginTime)
	*result.DeskId = t.DeskId
	for _, bean := range t.Records {
		b := bean.TransBeanUserRecord()
		result.Users = append(result.Users, b)
	}
	return result
}
