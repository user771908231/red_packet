package api

type MjRoom interface {
	CreateDesk(config interface{}) (MjDesk, error)
	GetDesk() MjDesk
	EnterUser(userId uint32, key string) error
}
