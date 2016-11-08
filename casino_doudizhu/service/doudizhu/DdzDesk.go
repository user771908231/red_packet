package doudizhu

//斗地主的desk
type DdzDesk struct {
	key string
}

//New一个Desk
func NewDdzDesk() *DdzDesk {
	return nil
}



//斗地主的桌子
//把数据同步到redis中去
func (d *DdzDesk) Update2Redis() error {
	return nil
}

//添加一个玩家
func (d *DdzDesk) AddUser() error {
	return nil
}


