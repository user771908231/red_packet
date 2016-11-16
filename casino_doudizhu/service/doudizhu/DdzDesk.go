package doudizhu

import (
	"errors"
	"casino_server/common/log"
	"sync"
	"casino_server/common/Error"
	"fmt"
)

//斗地主的desk
type DdzDesk struct {
	sync.Mutex
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

//得到一个用户
func (d *DdzDesk) GetUserByUserId(userId uint32) *DdzUser {
	for _, u := range d.Users {
		if u != nil && u.GetUserId() == userId {
			return u
		}
	}

	//哈哈哈
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

//都准备
func (d *DdzDesk) IsAllReady() bool {
	for _, user := range d.Users {
		if user != nil && user.IsNotReady() {
			return false
		}
	}

	return true
}

func (d *DdzDesk) CheckOutPai(out *POutPokerPais) error {
	right, err := out.GT(d.OutPai)
	if err != nil {
		log.E("出牌的时候，判断牌型的时候失败...")
		return Error.NewError(-1, "比较失败,出牌失败...")
	}

	if right {
		return nil
	} else {
		return Error.NewError(-1, "出的牌比别人的牌小，没有办法出牌")
	}

}

func (d *DdzDesk) GetUserIndexByUserId(userId uint32) int {
	for i, user := range d.Users {
		if user != nil && user.GetUserId() == userId {
			return i
		}
	}
	return -1;
}

func (d *DdzDesk) SetActiveUser(userId uint32) {
	*d.ActiveUser = userId
}

func (d *DdzDesk) SetDizhu(userId uint32) {
	*d.Dizhu = userId
}

//判断是否是当前活动玩家
func (d *DdzDesk) CheckActiveUser(userId uint32) error {
	if d.GetActiveUser() == userId {
		return nil
	} else {
		return Error.NewFailError(fmt.Sprintf("当前活动玩家是[%v]", d.GetActiveUser()))
	}
}

