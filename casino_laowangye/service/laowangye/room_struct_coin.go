package laowangye

import (
	"errors"
)

//查找一个空闲房间，如果没有空闲的则创建一个新房间
func (room *Room) GetFreeCoinDesk() (*Desk, error) {
	if room.GetRoomId() == 0 {
		return nil, errors.New("参数非法！")
	}
	for _,desk := range room.Desks{
		if desk == nil || desk.DeskOption == nil {
			continue
		}

		if len(desk.Users) < int(desk.DeskOption.GetMaxUser()) {
			return desk, nil
		}
	}

	//未找到有空闲位置的房间，则创建一个新房间
	new_desk, err := room.NewCoinDesk()
	if new_desk != nil && err == nil {
		//创建一个新房间
		room.Desks = append(room.Desks, new_desk)
		return new_desk, nil
	}

	return nil, errors.New("未找到空闲房间！")
}

//创建一个金币桌
func (room *Room) NewCoinDesk() (*Desk, error) {

	return nil, nil
}
