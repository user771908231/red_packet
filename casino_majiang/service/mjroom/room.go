package mjroom

//room 的接口设计
type Room interface {
	run()
	stop()
	createDesk()
	getDesk(key string)
}


//room的骨架设计
type RoomSkeleton struct {
	desks map[int32]*Desk
}


//在room中增加一个桌子
func (r *RoomSkeleton) addDesk(d Desk) {
	r.desks[d.GetDeskId()] = d
}
