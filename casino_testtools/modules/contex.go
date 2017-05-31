package modules

import (
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

//自定义Context模块
type Context struct {
	*macaron.Context
	Session session.Store
}
