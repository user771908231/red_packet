package doudizhu

import (
	"errors"
	"casino_server/common/log"
)

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
func (d *DdzDesk) AddUser(userId uint32) error {
	user := NewDdzUser()
	err := d.AddUserBean(user)
	return err
}

func (d *DdzDesk) AddUserBean(user *DdzUser) error {
	for i := 0; len(d.Users); i++ {
		if d.Users[i] == nil {
			d.Users[i] = user
			return nil
		}
	}
	log.E("玩家[%v]加入desk[%v]失败，因为没有合适的座位.", user.GetUserId(), d.GetDeskId())
	return errors.New("加入失败，没有找到合适的座位...")
}

