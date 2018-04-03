package redpack

import (
	"casino_redpack/modules"
	"fmt"
	"casino_redpack/model/redModel"
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/log"
	"time"
	"encoding/json"
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
		"msg": "%s",
		"redInfo": {
			"id": 10008,
		}
	}
}`, res_code, res_msg)
		ctx.Write([]byte(data))
	}()
	user_info := ctx.IsLogin()
	req_type := redModel.RoomType(ctx.QueryInt("type"))
	req_money := float64(ctx.QueryInt("money"))

	if user_info.Coin < req_money {
		return
	}
	//发红包五人
	room := redModel.GetRoomByType(req_type)
	room.SendRedpack(user_info, req_money, 5, 0)
	//减去用户的金币
	GetUserUplate(user_info,req_money,0)

	res_code = 1
	res_msg = fmt.Sprintf("发红包成功")


}

//五人对战：加入红包对战
func JoinWurenRedPacketHandler(ctx *modules.Context) {
	redId := ctx.QueryInt("redId")
	user := ctx.IsLogin()
	lists := redModel.GetRoomByType(redModel.RoomTypeWurenDZ).GetRedpackById(int32(redId))
	res := bson.M{
		"code": 0,
		"message": "faid",
		"request":bson.M{},
	}
	if lists != nil {
		lists.OpenRecord = append(lists.OpenRecord,&redModel.OpenRecordItem{
			UserId :user.Id, //领红包的人id
			NickName:user.HeadUrl,
		})

		res["code"] = 1
		res["message"] = "success"
		res["request"] = bson.M{
			"msg": "加入成功！",
			"redInfo": bson.M{
				"id":lists.Id,
				"nickname": lists.CreatorName,
				"headimgurl": lists.CreatorHead,
				"money": lists.Money,
				"all_membey": lists.Piece,
				"has_member": len(lists.OpenRecord) + 1,
			},
			"redItemList": []bson.M{
				bson.M{
					"headimgurl": user.HeadUrl,
				},
			},
		}
	}
	json_str,_ := ctx.JSONString(res)
	ctx.Write([]byte(json_str))
}

//炸弹接龙发红包
func SendZhadanRedPacketHandler(ctx *modules.Context) {
	res_code := 0
	res_msg := "金币不足，最低需要14金币。"

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
	//检查用户金币数
	if user_info.Coin < float64(14) {
		return
	}

	req_type := ctx.QueryInt("type")
	req_money := float64(ctx.QueryInt("money"))
	req_tailNumber := ctx.QueryInt("tailNumber")
	rep_number := ctx.QueryInt("nuber")
	if rep_number < 7 {
		rep_number = 7
	}
	switch req_type {
	case 5:
		//炸弹接龙（扫雷）
		room := redModel.GetRoomByType(redModel.RoomTypeSaoLei)
		if room == nil {
			res_msg = "该房间不存在，发红包失败！"
			return
		}

		_,err := room.SendRedpack(user_info, req_money, rep_number, req_tailNumber)
		if err != nil {
			res_msg = err.Error()
			return
		}
		//减去用户的金币

	}
	err := GetUserUplate(user_info,req_money,0)
	fmt.Println(err)
	res_code = 1
	res_msg = fmt.Sprintf("发红包成功！")
}

//扫雷，开红包按钮
func SaoleiJLOpenRedButtonAjaxHandler(ctx *modules.Context) {
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

	res_code = 1
	res_msg = "开红包成功！"

}

//扫雷 开红包记录
func SaoleiRedOpenRecordAjaxHandler(ctx *modules.Context) {
	res_code := 0

	res := bson.M{
		"code": &res_code,
		"message": "success",
		"request": bson.M{
			"redInfo": bson.M{},
			"user": bson.M{},
			"redItemList": []bson.M{},
		},
	}

	defer func() {
		res_json,_ := json.Marshal(res)
		ctx.Write(res_json)
	}()


	red_id := int32(ctx.QueryInt("redId"))
	red_info := redModel.GetRoomByType(redModel.RoomTypeSaoLei).GetRedpackById(red_id)

	if red_info == nil {
		log.E("SaoLei红包不存在！ id %d", red_id)
		return
	}

	red_status := 0
	if len(red_info.OpenRecord) == red_info.Piece {
		red_status = 1
	}
	//红包信息
	res["request"].(bson.M)["redInfo"] = bson.M{
		"id": red_info.Id,
		"type": red_info.Type,
		"status": red_status,
		"money": red_info.Money,
		"moneyfa": red_info.Money,
		"all_membey": red_info.Piece,
		"has_member": len(red_info.OpenRecord),
		"nickname": red_info.CreatorName,
		"headimgurl": red_info.CreatorHead,
		"tail_number": red_info.TailNumber,
	}

	user_info := ctx.IsLogin()

	//开红包
	open_money := red_info.Open(user_info)
	open_tail_num := int(open_money * 100)%10
	//开红包的玩家信息
	res["request"].(bson.M)["user"] = bson.M{
		"id": user_info.Id,
		"red_id": red_info.Id,
		"member_id": 1406,
		"money": open_money,
		"open_time": time.Now().Unix(),
		"open_status": 1,
		"winning": 1,
		"deduct_money": "0.01",
		"if_banker": 0,
		"join_money": "7.00",
		"win_money": "0.41",
		"primordial_money": "0.42",
		"code": open_tail_num,   //用户开的尾号
		"banker_name": "",
		"is_robot": 0,
		"tzok": 0,
		"nickname": user_info.NickName,
		"headimgurl": user_info.HeadUrl,
	}

	//开红包记录
	recore_list := []bson.M{}
	for _,item := range red_info.OpenRecord{
		new_item := bson.M{
			//"id": 9396,
			//"red_id": 68171,
			//"member_id": 667478,
			"money": item.Money,
			"open_time": item.Time.Unix(),
			//"open_status": 1,
			//"winning": 1,
			"deduct_money": item.Money * 0.03,  //扣除的钱
			//"if_banker": 0,
			//"join_money": "14.00",
			//"win_money": "1.21",
			//"primordial_money": "1.25",
			//"code": "5",
			//"banker_name": "",
			//"is_robot": 1,
			//"tzok": 0,
			"nickname": item.NickName,
			"headimgurl": item.Head,
		}
		recore_list = append(recore_list, new_item)
	}

	res["request"].(bson.M)["redItemList"] = recore_list


	res_code = 1
}

func GetRedPacketInfoHandler(ctx *modules.Context) {
	redId := ctx.QueryInt("redId")
	fmt.Println(redId)
	val := redModel.GetRoomByType(redModel.RoomTypeWurenDZ).GetRedpackById(int32(redId))
	res := bson.M{
		"code": 0,
		"message": "faid",
		"request":bson.M{},
	}
	if val != nil {
		res["code"] = 1
		res["message"] = "success"
		res["request"] = bson.M{
			"redInfo": bson.M{
				"id":val.Id,
				"nickname": val.CreatorName,
				"headimgurl": val.CreatorHead,
				"Money": val.Money,
				"all_membey": val.Piece,
				"has_member": len(val.OpenRecord) ,
			},
		}
	}

	json_str,_ := ctx.JSONString(res)
	ctx.Write([]byte(json_str))
}