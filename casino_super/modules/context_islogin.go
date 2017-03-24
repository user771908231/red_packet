package modules

import (
	"casino_super/model/userModel"
	"github.com/chanxuehong/wechat.v2/mp/oauth2"
)

//是否登录
func (ctx *Context) IsLogin() *userModel.User {
	user := ctx.Session.Get("user")
	if user == nil {
		return nil
	}else {
		user_info := user.(userModel.User)
		return &user_info
	}
}

//是否微信登录
func (ctx *Context) IsWxLogin() *oauth2.UserInfo {
	user := ctx.Session.Get("wx_user")
	var ret *oauth2.UserInfo = nil
	if user != nil {
		user_info := user.(oauth2.UserInfo)
		ret = &user_info
	}
	//ret = &oauth2.UserInfo{
	//	OpenId: "oab_twdGywSq99suV8Mt14FNxxOQ", //sapmm: oab_twdGywSq99suV8Mt14FNxxOQ wsx: oG9kZwrv8oFF9ja6WRlHxxMoJZoU 王沛：oab_twV764cO7tAl8IgCHlWgG2uw
	//	Nickname: "董兵",
	//	Sex: 1,      // 用户的性别, 值为1时是男性, 值为2时是女性, 值为0时是未知
	//	City: "成都",     // 普通用户个人资料填写的城市
	//	Province: "四川", // 用户个人资料填写的省份
	//	Country: "中国", // 国家, 如中国为CN
	//	// 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），
	//	// 用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	//	HeadImageURL: "http://wx.qlogo.cn/mmopen/ajNVdqHZLLDR9YkFYEz0XhumSbNtrpn98PlbDp7K87CxAGYMhkRwV6LEiaYPNRftBoktV2yXTQlodYEUA7SpZkg/0",
	//
	//	Privilege: []string{}, // 用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	//	UnionId: "oKKIfxHr0Gf-36iEgKDLCJUzeqrg",   // 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。 sapmm:oKKIfxHr0Gf-36iEgKDLCJUzeqrg   wsx:oKKIfxCxJiKHLOZcOobsBzb2sl0Q 王沛：oKKIfxDEIIkZVftBs9F1Yn8hzMCg
	//}
	return ret
}
