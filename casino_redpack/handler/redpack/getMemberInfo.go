package redpack

import (
	"casino_redpack/modules"
	"fmt"
	"casino_redpack/model/userModel"
	"errors"
	"casino_redpack/model/redModel"
	"casino_common/common/log"
)
//获取用户信息
func GetMemberInfo(ctx *modules.Context) {
	res := `{
	"code": 0,
	"message": "error"}`
	User := userModel.GetUserById(ctx.IsLogin().Id)
	if User == nil {
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

	res = fmt.Sprintf(res, User.Id,User.NickName, User.NickName, User.HeadUrl, User.Coin)
	ctx.Write([]byte(res))
	return
}


func GetUserUplate(user *userModel.User,money float64,Type int) error {
	XinUser := userModel.GetUserById(user.Id)
	if Type == 0 {

		err := XinUser.CapitalUplete("-",money)
		if err != nil {
			return errors.New("减去用户金币失败！")
		}
		data := redModel.CoinAddSbtract{
			UserId:user.Id,
			SendOrOpenPacket:Type,
			UserCoin:user.Coin,
			AddOrSubtract:money,
		}
		err1 := data.Isert()
		if err1 != nil {
			XinUser.CapitalUplete("+",money)
			log.E("用户【%s】减去金币时生成记录失败！",user.NickName)
			return err1
		}
		return nil
	}else {
		err := XinUser.CapitalUplete("+",money)
		if err != nil {
			return errors.New("减去用户金币失败！")
		}
		data := redModel.CoinAddSbtract{
			UserId:user.Id,
			SendOrOpenPacket:Type,
			UserCoin:user.Coin,
			AddOrSubtract:money,
		}
		err1 := data.Isert()
		if err1 != nil {
			user.CapitalUplete("-",money)
			log.E("%s",err1)
			return err1
		}

		return nil
	}

}


