package autoLogin

import (
	"casino_redpack/modules"
	"casino_common/common/log"
	"casino_redpack/model/autoLoginModel"


	"time"
	"casino_common/common/userService"
	"casino_redpack/model/userModel"
	"strings"
)

//接受传来的数据
func AcceptData(ctx *modules.Context) {
	data := strings.Replace(ctx.Query("data"), " ", "+", -1)
	log.T("接收到Data的值：%s",data)
	key := []byte("123asdssssssssss")
	//解密
	str,err := autoLoginModel.DataDecoding(data,key)
	if err == nil {
		string := string(str)
		log.T("转化为string的值：%s",string)
		data := autoLoginModel.CharacterSplit(string)
		if data != nil {
			// Sub 计算两个时间差
			subM := time.Now().Sub(data.Time)
			if subM.Minutes() <  float64(10) {
				//获取用户信息
				user_info := userService.GetUserById(data.Id)
				//根据ID查找本游戏是否有相同
				if user_info != nil {
					user1 := userModel.GetUserByThreePartyId(*user_info.Id)
					//每次登陆更新获取到的三方用户信息
					user1.Coin = float64(*user_info.Coin)
					user1.NickName = *user_info.NickName
					user1.HeadUrl = *user_info.HeadUrl
					user1.UnionId = *user_info.UnionId
					user1.OpenId = *user_info.OpenId
					user1.PassWd = *user_info.Pwd
					if user1 == nil {
						user := userModel.User{
							ThreePartyId:*user_info.Id,
							NickName:*user_info.NickName,
							HeadUrl:*user_info.HeadUrl,
							OpenId:*user_info.OpenId,
							UnionId:*user_info.UnionId,
							PassWd:*user_info.Pwd,
							Coin:float64(*user_info.Coin),
						}
						err := user.Insert()
						if err == nil {
							log.T("注册信息成功！")
							user2 := userModel.GetUserByThreePartyId(*user_info.Id)
							ctx.Session.Set("user", *user2)
							ctx.Redirect("/home",302)
							log.T("用户[%d]信息获取成功！写入Session，跳转home",user_info.Id)
							return
						}
						log.T("注册信息失败！错误信息：",err)
						return
					}
					//更新用户信息
					err := user1.Uplate()
					if err != nil {
						log.T("更新信息失败")
					}
					log.T("用户的金币",*user_info.Coin)
					//获取更新h
					ctx.Session.Set("user", *user1)
					ctx.Redirect("/home",302)
					log.T("用户[%d]信息获取成功！写入Session，跳转home",user_info.Id)
					return

				}
				log.T("用户信息获取失败")
				return
			}else {
				log.T("Id:%d 登陆超时",data.Id)
				return
			}
			log.T("this is data nil")
			return
		}
	}
	log.T("解密失败！")
	return
}
