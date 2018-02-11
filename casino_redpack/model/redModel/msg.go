package redModel

import (
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
)

//解析成客户端需要的红包列表
func GetClientRedpackListJson(list []*Redpack) []byte {
	res := bson.M{
		"code": 1,
		"message": "success",
		"request": []bson.M{},
	}

	res_list := []bson.M{}
	for _,item := range list {
		res_list = append(res_list, bson.M{
			"id": item.Id,
			"type": item.Type,
			"money": item.Money,
			"moneyfa": item.Lost,
			"member_id": item.CreatorUser,
			"all_membey": 10,
			"has_member": 9,
			"tail_number": item.TailNumber,
			"nickname": item.CreatorName,
			"headimgurl": item.CreatorHead,
			"xianshi": 0,
			"jiaru": 0,
		})
	}

	res["request"] = res_list
	data,_ := json.Marshal(res)
	return data
}
