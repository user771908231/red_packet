package data

import "casino_majiang/service/majiang"

//麻将桌子的一些状态
type MjUserStatus struct {
	Status        int32
	Ready         bool
	DingQue       bool
	Exchange      bool
	AgentMode     bool
	IsBanker      bool
	IsBreak       bool
	IsLeave       bool
	IsRobot       bool
	ApplyDissolve int32 //是否同意解散房间
}

//todo
func (s *MjUserStatus) IsReady() bool {
	return false
}

//是否游戏中
func (s *MjUserStatus) IsGaming() bool {
	if s.Status == majiang.MJDESK_STATUS_RUNNING {
		return true
	} else {
		return false
	}
}

//用户是否胡牌
func (s *MjUserStatus) IsHu() bool {
	return s.Status == majiang.MJUSER_STATUS_HUPAI
}

//用户是否未胡牌
func (s *MjUserStatus) IsNotHu() bool {
	return !s.IsHu()
}

//设置用户的状态
func (s *MjUserStatus) SetStatus(status int32) {
	s.Status = status
}

func (s *MjUserStatus) GetAgentMode() bool {
	return s.AgentMode
}
func (s *MjUserStatus) SetAgentMode(a bool) {
	s.AgentMode = a
}

//
func (s *MjUserStatus) GetApplyDissolve() int32 {
	return s.ApplyDissolve
}
