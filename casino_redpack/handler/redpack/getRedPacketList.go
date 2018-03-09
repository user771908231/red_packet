package redpack

import (
	"casino_redpack/modules"
	"gopkg.in/mgo.v2/bson"
	"casino_redpack/model/redModel"
)

//五人对战->房间列表
func GetRedPacketListHandler(ctx *modules.Context) {
	res := bson.M{
		"code": 1,
		"message": "success",
		"request": []bson.M{},
	}

	req_type := ctx.QueryInt("type")

	list := []bson.M{}

	room := redModel.GetRoomByType(redModel.RoomType(req_type))
	for _,item := range room.RedpackList {
		list = append(list, bson.M{
			"id": item.Id,
			"type": item.Type,
			"money": item.Money,
			"all_membey": item.Piece,
			"has_member": item.Piece - len(item.OpenRecord),
			"tail_number": item.TailNumber,
			"nickname": item.CreatorName,
			"headimgurl": item.CreatorHead,
		})
	}

	res["request"] = list
	json_res,_ := ctx.JSONString(res)
	ctx.Write([]byte(json_res))
}
