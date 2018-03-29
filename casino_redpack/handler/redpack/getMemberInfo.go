package redpack

import (
	"casino_redpack/modules"
	"fmt"
)

//获取用户信息
func GetMemberInfo(ctx *modules.Context) {
	res := `{
	"code": 0,
	"message": "error"}`

	user_info := ctx.IsLogin()
	if user_info == nil {
		ctx.Write([]byte(res))
		return
	}

	res = `{
	"code": 1,
	"message": "success",
	"request": {
		"id": %d,
		"username": "%s",
		"nickname": "%s",
		"headimgurl": "%s",
		"golds": "%.2f",
		"is_agent": 1,
		"numberOnline": 1238
	}
}`

	res = fmt.Sprintf(res, user_info.Id, user_info.NickName, user_info.NickName, user_info.HeadUrl, user_info.Capital)
	ctx.Write([]byte(res))
	return
}
