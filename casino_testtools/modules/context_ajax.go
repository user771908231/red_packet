package modules

import "gopkg.in/mgo.v2/bson"

//ajax方法
func (ctx *Context) Ajax(code int, msg string, data interface{}) {
	ctx.JSON(200, bson.M{
		"code": code,
		"msg":  msg,
		"data": data,
	})
	ctx.Resp.Write([]byte{})
}
