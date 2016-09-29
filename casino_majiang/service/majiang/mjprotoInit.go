package majiang

func NewMjpai() *MJPai {
	ret := &MJPai{}
	ret.Value = new(int32)
	ret.Des = new(string)
	ret.Flower = new(int32)
	ret.Index = new(int32)
	return ret
}

//返回一个麻将room
func NewMjRoom() *MjRoom {
	ret := &MjRoom{}
	return ret
}

//返回一个麻将
func NewMjDesk() *MjDesk {
	ret := &MjDesk{}
	ret.DeskId = new(int32)
	ret.RoomId = new(int32)
	ret.Status = new(int32)
	ret.Password = new(string)
	ret.Owner = new(uint32)
	ret.CreateFee = new(int64)
	ret.MjRoomType = new(int32)
	ret.BoardsCout = new(int32)
	ret.CapMax = new(int64)
	ret.CardsNum = new(int32)
	ret.Settlement = new(int32)
	ret.BaseValue = new(int64)
	ret.ZiMoRadio = new(int32)
	ret.DianGangHuaRadio = new(int32)
	return ret
}