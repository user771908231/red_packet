package redModel

import (
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
)

//解析成客户端需要的红包列表
func GetClientRedpackListJson(list []*Redpack, userId uint32) []byte {
	res := bson.M{
		"code": 1,
		"message": "success",
		"request": []bson.M{},
	}

	res_list := []bson.M{}
	//lenth := len(list)
	for _,item := range list {
		//只取最后十条
		//if i < lenth-10 {
		//	continue
		//}
		renshu := len(item.OpenRecord)
		if renshu == item.Piece {
			continue
		}
		new_item := bson.M{
			"id": item.Id,
			"type": item.Type,
			"money": item.Money,
			"moneyfa": item.Money,
			"member_id": item.CreatorUser,
			"all_membey": 10,
			"has_member": 9,
			"tail_number": item.TailNumber,
			"nickname": item.CreatorName,
			"headimgurl": item.CreatorHead,
			"is_self": 0,
			"jiaru": 0,
		}

		//是否本人发的
		if userId == item.CreatorUser {
			new_item["is_self"] = 1
		}
		//是否已领取
		for _,u := range item.OpenRecord {
			if userId == u.UserId {
				new_item["jiaru"] = 1
			}
		}
		if renshu == item.Piece {
			new_item["jiaru"] = 2
		}

		res_list = append(res_list, new_item)
	}

	res["request"] = res_list
	data,_ := json.Marshal(res)
	return data
}
