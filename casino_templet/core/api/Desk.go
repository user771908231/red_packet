package api

type DeskApi interface {
	GetDeskId() int32                  //得到Desk的Id
	GetDeskPassword() string           //得到desk的房间号
	GetDeskConfig() interface{}        //得到config
	GetDeskSkeleton() interface{}      //得到骨架
	SetRoom(interface{}) error         //设置房间
	ActEnterUser(...interface{}) error //进入房间
	ActReady(...interface{}) error     //准备
}

type DeskCore struct {
}

func (*DeskCore) GetDeskId() int32 {
	panic("implement me")
}

func (*DeskCore) GetDeskPassword() string {
	panic("implement me")
}

func (*DeskCore) GetDeskConfig() interface{} {
	panic("implement me")
}

func (*DeskCore) GetDeskSkeleton() interface{} {
	panic("implement me")
}

func (*DeskCore) SetRoom(interface{}) error {
	panic("implement me")
}

func (*DeskCore) ActEnterUser(...interface{}) error {
	panic("implement me")
}

func (*DeskCore) ActReady(...interface{}) error {
	panic("implement me")
}
