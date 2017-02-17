package data

//麻将桌子的一些状态
type MjDeskStatus struct {
	status int32
}

//todo
func (s *MjDeskStatus) IsNotGaming() bool {
	return false
}

func (s *MjDeskStatus) IsGaming() bool {
	return false
}

//todo
func (s *MjDeskStatus) IsNotPreparing() bool {
	return false
}

func (s *MjDeskStatus) S() int32 {
	return s.status
}

//设置status
func (s *MjDeskStatus) SetStatus(ss int32) {
	s.status = ss
}
