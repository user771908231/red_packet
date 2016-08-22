package mode

import "time"



//一个人一场游戏
type BeanRecord struct {
	UserId    uint32
	NickName  string
	WinAmount int64
}


//一局游戏的数据
type T_th_desk_record struct {
	DeskId     int32
	GameNumber int32
	UserIds    string
	BeginTime  time.Time
	Records    []BeanRecord
}


