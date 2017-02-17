package data

import "casino_majiang/service/majiang"

//麻将桌子的一些状态
type MjDeskStatus struct {
	Status     int32
	IsExchange bool
}

func (s *MjDeskStatus) IsNotGaming() bool {
	return !s.IsNotGaming()
}

func (s *MjDeskStatus) IsNotPreparing() bool {
	return !s.IsPreparing()
}

func (s *MjDeskStatus) S() int32 {
	return s.Status
}

//设置status
func (s *MjDeskStatus) SetStatus(ss int32) {
	s.Status = ss
}

//是否已经开始游戏了...
func (s *MjDeskStatus) IsGaming() bool {
	if s.S() == majiang.MJDESK_STATUS_RUNNING {
		return true
	} else {
		return false
	}
}

//这里表示 是否是 [正在] 准备中...
func (s *MjDeskStatus) IsPreparing() bool {
	if s.S() == majiang.MJDESK_STATUS_READY {
		return true
	} else {
		return false
	}
}
