package majiang


//普通的麻将房间...

//room的结构定义在proto中
//type MjRoom struct {
//	roomType int32
//
//}

func (r *MjRoom) CreateDesk() *MjDesk {
	//create 的时候，是否需要通过type 来判断,怎么样创建房间
	return nil
}

//room中增加一个desk
func (r *MjRoom) AddDesk() error {
	return nil

}

