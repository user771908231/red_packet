package data

import "casino_majiang/service/majiang"

//麻将桌子的一些状态
type MjUserStatus struct {
	status int32
	Ready  bool
}

//todo
func (s *MjUserStatus) IsReady() bool {
	return false
}

//是否游戏中
func (s *MjUserStatus) IsGaming() bool {
	if s == majiang.MJDESK_STATUS_RUNNING {
		return true
	} else {
		return false
	}
}


//用户是否胡牌
func (s *MjUserStatus) IsHu() bool {
	return s == majiang.MJUSER_STATUS_HUPAI
}

//用户是否未胡牌
func (s *MjUserStatus) IsNotHu() bool {
	return !s.IsHu()
}

//设置用户的状态
func (s *MjUserStatus) SetStatus(status int32) {
	s.status = status
}
