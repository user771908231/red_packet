package redpack

import (
	"casino_redpack/modules"
	"fmt"
	"casino_redpack/model/redModel"
	"gopkg.in/mgo.v2/bson"
)

//五人对战：发红包
func SendWurenRedPacketHandler(ctx *modules.Context) {
	res_code := 0
	res_msg := "金币不足，最低需要10金币。"

	defer func() {
		data := fmt.Sprintf(`{
	"code": 1,
	"message": "success",
	"request": {
		"code": %d,
		"msg": "%s"
	}
}`, res_code, res_msg)
		ctx.Write([]byte(data))
	}()

	req_type := ctx.QueryInt("type")
	req_money := ctx.QueryInt("money")

	switch req_type {

	}

	res_code = 0
	res_msg = fmt.Sprintf("金币不足，最低需要%d金币。", req_money)
}

//五人对战：加入红包对战
func JoinWurenRedPacketHandler(ctx *modules.Context) {
	res := bson.M{
		"code": 1,
		"message": "success",
		"request": bson.M{
			"msg": "加入成功！",
			"redInfo": bson.M{
				"nickname": "郑细弟",
				"headimgurl": "http://wx.qlogo.cn/mmopen/ajNVdqHZLLDR9YkFYEz0XhumSbNtrpn98PlbDp7K87CxAGYMhkRwV6LEiaYPNRftBoktV2yXTQlodYEUA7SpZkg/0",
				"money": 10,
				"all_membey": 5,
				"has_member": 5,
			},
			"redItemList": []bson.M{
				bson.M{
					"headimgurl": "http://wx.qlogo.cn/mmopen/ajNVdqHZLLDR9YkFYEz0XhumSbNtrpn98PlbDp7K87CxAGYMhkRwV6LEiaYPNRftBoktV2yXTQlodYEUA7SpZkg/0",
				},
			},
		},
	}

	json_str,_ := ctx.JSONString(res)
	ctx.Write([]byte(json_str))
}

//炸弹接龙发红包
func SendZhadanRedPacketHandler(ctx *modules.Context) {
	res_code := 0
	res_msg := "金币不足，最低需要10金币。"

	defer func() {
		data := fmt.Sprintf(`{
	"code": 1,
	"message": "success",
	"request": {
		"code": %d,
		"msg": "%s"
	}
}`, res_code, res_msg)
		ctx.Write([]byte(data))
	}()

	//------------------------------------------------
	user_info := ctx.IsLogin()
	if user_info == nil {
		res_msg = "清先登录！"
		return
	}

	req_type := ctx.QueryInt("type")
	req_money := float64(ctx.QueryInt("money"))
	req_tailNumber := ctx.QueryInt("tailNumber")

	switch req_type {
	case 5:
		room := redModel.GetRoomById(10000)
		if room == nil {
			res_msg = "该房间不存在，发红包失败！"
			return
		}

		_,err := room.SendRedpack(user_info, req_money, 10, req_tailNumber)
		if err != nil {
			res_msg = err.Error()
			return
		}
	}

	res_code = 1
	res_msg = fmt.Sprintf("发红包成功！")
}
