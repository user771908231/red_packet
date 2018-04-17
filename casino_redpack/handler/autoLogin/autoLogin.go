package autoLogin

import (
	"casino_redpack/modules"
	"casino_common/common/log"
	"casino_redpack/model/autoLoginModel"


	"time"
	"casino_common/common/userService"
)

//接受传来的数据
func AcceptData(ctx *modules.Context) {
	data := ctx.Query("data")
	log.T("接收到Data的值：%s",data)
	key := []byte("")
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
					ctx.Session.Set("user", *user_info)
					ctx.Redirect("/home",302)
					log.T("用户[%d]信息获取成功！写入Session，跳转home",user_info.Id)
				}
				log.T("用户信息获取失败")
			}else {
				log.T("Id:%d 登陆超时",data.Id)
			}
			log.T("this is data nil")
		}
	}
	log.T("解密失败！")
}
