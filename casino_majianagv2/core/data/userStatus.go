package data

//麻将桌子的一些状态
type MjUserStatus struct {
	status int32
	Ready  bool
}

//todo
func (s *MjUserStatus) IsReady() bool {
	return false
}

//设置用户的状态
func (s *MjUserStatus) SetStatus(status int32) {
	s.status = status
}
