package modules

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
)

//自定义Context模块
type Context struct {
	*macaron.Context
	Session session.Store
}