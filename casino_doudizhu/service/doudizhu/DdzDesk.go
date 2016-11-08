package doudizhu

//斗地主的desk
type DdzDesk struct {
	*PDdzDesk
	Users []*DdzUser
}

//斗地主的桌子//把数据同步到redis中去
func (d *DdzDesk) Update2Redis() error {
	bak := NewPDdzbak()
	bak.Desk = d.PDdzDesk
	for _, u := range d.Users {
		bak.Users = append(d.Users, u)
	}

	//备份desk的数据
	UpdateDesk2Redis(bak)
	return nil
}

//添加一个玩家
func (d *DdzDesk) AddUser() error {
	return nil
}





