package api

type PDKRoom interface {
	CreateDesk(interface{}) (PDKDesk, error) //创建房间
}
