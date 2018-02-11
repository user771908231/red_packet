package redpack

import (
	"casino_redpack/modules"
	"gopkg.in/mgo.v2/bson"
)

//五人对战->房间列表
func GetRedPacketListHandler(ctx *modules.Context) {
	res := bson.M{
		"code": 1,
		"message": "success",
		"request": []bson.M{},
	}

	list := []bson.M{
		bson.M{
			"id": 25287,
			"type": 1,
			"money": "10.00",
			"all_membey": 5,
			"has_member": 5,
			"tail_number": 0,
			"nickname": "郑细弟",
			"headimgurl": "http://wx.qlogo.cn/mmopen/ajNVdqHZLLDR9YkFYEz0XhumSbNtrpn98PlbDp7K87CxAGYMhkRwV6LEiaYPNRftBoktV2yXTQlodYEUA7SpZkg/0",
		},
	}

	res["request"] = list
	json_res,_ := ctx.JSONString(res)
	ctx.Write([]byte(json_res))
}
