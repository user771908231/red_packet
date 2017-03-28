package majiang

type MjuserChangShaConfig struct {
	changshaGang bool //长沙杠
}

func (u *MjuserChangShaConfig) GetChangShaGangStatus() bool {
	if u == nil {
		return false
	} else {
		return u.changshaGang
	}
}
