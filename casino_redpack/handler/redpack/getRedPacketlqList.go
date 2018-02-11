package redpack

import "casino_redpack/modules"

//红包列表
func GetRedPacketlqList(ctx *modules.Context) {
	list := `{
	"code": 1,
	"message": "success",
	"request": "0"
}`
	ctx.Write([]byte(list))
}


//

