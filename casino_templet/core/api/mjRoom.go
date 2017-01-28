package api

type MjRoom interface {
	CreateDesk(config interface{}) (error, MjDesk)
	GetDesk() MjDesk
	EnterUser(userId uint32, key string) error
}
