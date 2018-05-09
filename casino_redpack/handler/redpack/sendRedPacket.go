package redpack

import (
	"casino_redpack/modules"
	"fmt"
	"casino_redpack/model/redModel"
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/log"
	"time"
	"encoding/json"
	"casino_redpack/model/userModel"
	"math"
	"errors"
	"strings"
	"unicode/utf8"
	"regexp"
	"strconv"
	"casino_redpack/model/agentModel"
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

	if user_info.Coin < req_money * 1.7 {
		res_code = 0
		res_msg = fmt.Sprintf("发红包成功失败！最少需要有%.0f",req_money * 1.7)
		return
	}
	//发红包五人
	room := redModel.GetRoomByType(req_type)
	_,err:= room.SendRedpack(user_info, req_money, 5, 0)
	if err != nil {
		log.T("发红包错误",err)
	}
	//减去用户的金币
	user_info.CapitalUplete("-",req_money,"发红包")
	//GetUserUplate(user_info,req_money,0,"发红包")

	res_code = 1
	res_msg = fmt.Sprintf("发红包成功")


}

//五人对战：加入红包对战
func JoinWurenRedPacketHandler(ctx *modules.Context) {
	redId := ctx.QueryInt("redId")
	fmt.Println(redId)
	user := ctx.IsLogin()
	lists := redModel.GetRoomByType(redModel.RoomTypeWurenDZ).GetRedpackById(int32(redId))
	res := bson.M{
		"code": 0,
		"message": "faid",
		"request":bson.M{},
	}
	if lists != nil {
		//for _,item := range lists.OpenRecord{
		//	if item.UserId == user.Id {
		//		continue
		//	}else{
		//		lists.OpenRecord = append(lists.OpenRecord,&redModel.OpenRecordItem{
		//			UserId :user.Id, //领红包的人id
		//			NickName:user.HeadUrl,
		//		})
		//	}
		//}
		var redItemList []bson.M
		lists.OpenRecord = append(lists.OpenRecord,&redModel.OpenRecordItem{
			UserId :user.Id, //领红包的人id
			NickName:user.HeadUrl,
		})
		redItemList = append(redItemList,bson.M{
			"headimgurl": user.HeadUrl,
		},)

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
				"has_member": len(lists.OpenRecord),
			},
			"redItemList": redItemList,
				//[]bson.M{
				//bson.M{
				//	"headimgurl": user.HeadUrl,
				//},
			//},
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
	user_info := ctx.CurrentUserInfo()
	if user_info == nil {
		res_msg = "清先登录！"
		return
	}


	req_type := ctx.QueryInt("type")
	req_money := float64(ctx.QueryInt("money"))
	req_tailNumber := ctx.QueryInt("tailNumber")
	rep_number := ctx.QueryInt("nuber")
	log.T("发红包类型：",req_type)
	log.T("发红包大小：",req_money)
	log.T("发红包雷号：",req_tailNumber)
	log.T("发红包几份：",rep_number)
	//检查用户金币数
	if ctx.CurrentUserInfo().Coin < (req_money * 1.4) {
		msg := fmt.Sprintf("金币不足，你最少需要有%d金币。",int(req_money * 1.4))
		res_msg = msg
		log.T(msg)
		return
	}
	if rep_number < 5 {
		rep_number = 5
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
	GetUserUplate(user_info,req_money,0,"红包炸弹发红包")
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
	//红包Id
	redId := ctx.QueryInt("redId")
	//房间类型
	Type := ctx.QueryInt("Type")
	var Types redModel.RoomType
	//判断放假类型
	switch Type {
	case 1:
		Types = redModel.RoomTypeWurenDZ
	case 2:
		Types = redModel.RoomTypeNiuniu
	case 4:
		Types = redModel.RoomTypeErBaGang
	case 5:
		Types = redModel.RoomTypeSaoLei
	}
	//根基房间类型和红包ID获取红包信息
	info := redModel.GetRoomByType(Types).GetRedpackById(int32(redId))
	//判断用户的金币是否小于红包的金币大小
	var piece float64
	//判断红包的份数 给出用户开红包需要的金币倍率
	switch info.Piece {
	case 5:
		piece = 2
	case 6:
		piece = 1.7
	case 7:
		piece = 1.5
	case 8:
		piece = 1.3
	case 9:
		piece = 1.2
	case 10:
		piece = 1.1
	}
	val := FloatValue(info.Money * piece,0)
	if ctx.CurrentUserInfo().Id != info.CreatorUser {
		if ctx.CurrentUserInfo().Coin < val {
			//扣除红包金额的倍数
			res_msg = fmt.Sprintf("金币不足，最低需要%d金币。",int(val))
			return
		}
	}


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
	bools,open_money := red_info.Open(user_info)
	//open_tail_num := int(open_money * 100)%10 //获得的值有误差
	open_tail_num := GetWeishu(open_money)
	if user_info.AccountNumber == "" {
		user_info.AccountNumber = user_info.NickName
	}
	//开红包的玩家信息
	res["request"].(bson.M)["user"] = bson.M{
		"id": user_info.Id,
		"red_id": red_info.Id,
		"member_id": 1406,
		"money": fmt.Sprintf("%.2f", open_money),
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
		"nickname": user_info.AccountNumber,
		"headimgurl": user_info.HeadUrl,
	}
	log.T("开包的尾数：%d",open_tail_num)
	//判断是否记录过
	if bools {
		//判读开包数字
		str,value := IsNumberType(open_money)
		//获取用户
		if value != 0 {
			user_info.CapitalUplete("+",value,str)
		}
		go RebateLog(user_info,open_money)

		//判断用户是否中雷
		go JudgeInMine(open_tail_num,red_info.TailNumber,red_info.Money,open_money,red_info.Piece,user_info.Id,red_info.CreatorUser)


	}

	//开红包记录
	recore_list := []bson.M{}
	for _,item := range red_info.OpenRecord{
		new_item := bson.M{
			//"id": 9396,
			//"red_id": 68171,
			//"member_id": 667478,
			"money": fmt.Sprintf("%.2f",item.Money),
			"open_time": item.Time.Unix(),
			//"open_status": 1,
			//"winning": 1,
			"deduct_money": FloatValue(item.Money * 0.03,2),  //扣除的钱
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
/*
open_tail_num 开包的尾数
tailnumber 发包人设置的尾数
red_money	红包大小
money		开了多少钱
number		几人包
Odds      	赔率
ThisUserID	当前开包人ID
SendPacketUserId 	发包人ID
 */
func JudgeInMine(open_tail_num int,tailnumber int,red_money float64,money float64,number int,ThisUserID uint32,SendPacketUserId uint32) error{
	var err error
	//判断几人包设置赔率
	//7人包赔率1.6倍，8人包赔率1.4倍，9人包1.2 10人包赔率1倍
	var Odds float64
	switch number {
	case 5:
		Odds = 2
	case 6:
		Odds = 1.7
	case 7:
		Odds = 1.5
	case 8:
		Odds = 1.3
	case 9:
		Odds = 1.2
	case 10:
		Odds = 1.1
	default:
		Odds = 1
	}

	//判断用户是否中雷
	if  open_tail_num == tailnumber {
		//--中雷---
		log.T("中雷用户ID:%d",ThisUserID)
		this_user := userModel.GetUserById(ThisUserID)
		if this_user == nil {
			log.T("没有找到用户")
		}
		//open_record_multiple = 0.03 //百分之三的扣费率
		//要给开包玩家的金币数
		money1 := money - FloatValue(money * 0.03,2)
		err = this_user.CapitalUplete("+",money1,"开红包")
		if err != nil {
			log.T("用户ID：%d 加金币失败",ThisUserID)
		}
		log.T("开包玩家的金币数：",money1)
		//要给的发包玩家的金币数
		log.T("红包大小：",red_money)
		// 红包是自己的就不赔
		if ThisUserID == SendPacketUserId {
			return nil
		}
		money0 := FloatValue(red_money * Odds,2)
		//因为要赔给发包玩家金币大于开包玩家得到的金币 不做金币加 减去开包玩家的差值金币
		log.T("中雷赔的金币数：",money0)
		err = this_user.CapitalUplete("-",money0,"输")
		if err != nil {
			log.T("用户ID：%d 减金币失败",ThisUserID)
		}
		////获取推广用户
		//level_user := userModel.GetUserById(uint32(this_user.ExtensionId))
		//if level_user != nil {
		//	go AgentRebate(level_user,money0)
		//}

		//赔给发红包的玩家
		//获取发包人的信息
		SendUser := userModel.GetUserById(SendPacketUserId)
		if SendUser == nil {
			log.T("没有找到发包用户")
		}
		err := SendUser.CapitalUplete("+",money0,"赢")
		if err != nil {
			log.T("发包用户ID：%d 加金币失败",SendUser.Id)
		}
		return err
		//---end
	}else{
		//--没中雷---
		//open_record_multiple = 0.03 //百分之三的扣费率
		this_user := userModel.GetUserById(ThisUserID)
		if this_user == nil {
			log.T("没有找到此用户！ID:%d",ThisUserID)
			return errors.New("没有找到此用户！",)
		}
		money0 := money - FloatValue(money * 0.03,2)
		err = this_user.CapitalUplete("+",money0,"开红包")
		//---end
		if err != nil {
			return err
			log.T("➕操作错误")
		}
		return nil
	}
	//扣费记录
	//。。。

}

func FloatValue(f float64,n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}
//判读开的红包数字类型 123-789 1234-6789  111-999 1111-9999 001 1314 2018 520
func IsNumberType(f float64) (string,float64) {
	log.T("得到的值%f",f)
	s :=fmt.Sprintf("%.2f", f)
	log.T("得到的值字符串",s)
	ss := strings.Split(s,".")
	log.T("得到的值字符串去除小数点byte",ss)
	sss := strings.Join(ss,"")
	log.T("得到的去除小数点字符串",sss)
	slengt := utf8.RuneCountInString(sss)
	log.T("得到的去除小数点字符串的长度",slengt)
	switch slengt {
	case 3:
		//顺子
		match1, _ := regexp.MatchString("^(?:(123)||(234)||(345)||(456)||(567)||(678)||(789)){3}$", sss)
		//豹子
		//match2, _ := regexp.MatchString("^(?:([1-9])\1{/\2})$", sss)
		//520
		match7, _ := regexp.MatchString("^(?:(520){3})$", sss)
		//001
		match9, _ := regexp.MatchString("^(?:(001)){3}$", sss)
		if match1 {
			log.T("3顺子",6.88)
			return "顺子",6.88
		}else if sss == "111" || sss == "222" || sss =="333" || sss == "444" || sss == "555" || sss == "666" || sss == "777" || sss == "888" || sss == "999"{
			log.T("3豹子",10)
			return "豹子",10
		}else if match7 {
			log.T("520",18.88)
			return "5.20",18.88
		}else if match9 {
			log.T("001",18.88)
			return "0.01",18.88
		}
		log.T("什么都不是",sss)
		return "",0
		break;
	case 4:
		match3, _ := regexp.MatchString("^(?:(1234)||(2345)||(3456)||(4567)||(5678)||(6789)){4}$", sss)

		//match4, _ := regexp.MatchString("^(?:([0-9])\1){2}$", sss)

		match8, _ := regexp.MatchString("^(?:(1314)){4}$", sss)

		match0, _ := regexp.MatchString("^(?:(2018)){4}$", sss)

		if match3 {
			log.T("4顺子",38.88)
			return "顺子",38.88
		}else if sss == "1111" || sss == "2222" || sss =="3333" || sss == "4444" || sss == "5555" || sss == "6666" || sss == "7777" || sss == "8888" || sss == "9999" {
			log.T("4豹子",58.88)
			return "豹子",58.88
		}else if match8 {
			log.T("1314",18.88)
			return "1314",18.88
		}else if match0 {
			log.T("2018",18.88)
			return "2018",18.88
		}
		log.T("什么都不是",sss)
		return "",0
		break;
	case 5:
		match5, _ := regexp.MatchString("^(?:(12345)||(23456)||(34567)||(45678)||(56789)){5}$", sss)

		//match6, _ := regexp.MatchString("^(?:([0-9])\1){4}$", sss)

		if match5 {
			log.T("5顺子",388)
			return "顺子",388
		}else if sss == "11111" || sss == "22222" || sss =="33333" || sss == "44444" || sss == "55555" || sss == "66666" || sss == "77777" || sss == "88888" || sss == "99999"{
			log.T("5豹子",588)
			return "豹子",588
		}
		log.T("什么都不是",sss)
		return "",0
		break;
	default:
		log.T("什么都不是",sss)
		return "",0
		break;
	}
	return "",0
}

//五人对战：加入红包对战
func JoinNunuRedPacketHandler(ctx *modules.Context) {
	redId := ctx.QueryInt("redId")
	fmt.Println(redId)
	user := ctx.IsLogin()
	lists := redModel.GetRoomByType(redModel.RoomTypeNiuniu).GetRedpackById(int32(redId))
	res := bson.M{
		"code": 0,
		"message": "faid",
		"request":bson.M{},
	}
	if lists != nil {
		//for _,item := range lists.OpenRecord{
		//	if item.UserId == user.Id {
		//		continue
		//	}else{
		//		lists.OpenRecord = append(lists.OpenRecord,&redModel.OpenRecordItem{
		//			UserId :user.Id, //领红包的人id
		//			NickName:user.HeadUrl,
		//		})
		//	}
		//}
		var redItemList []bson.M
		lists.OpenRecord = append(lists.OpenRecord,&redModel.OpenRecordItem{
			UserId :user.Id, //领红包的人id
			NickName:user.HeadUrl,
		})
		redItemList = append(redItemList,bson.M{
			"headimgurl": user.HeadUrl,
		},)

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
				"has_member": len(lists.OpenRecord),
			},
			"redItemList": redItemList,
			//[]bson.M{
			//bson.M{
			//	"headimgurl": user.HeadUrl,
			//},
			//},
		}
	}
	json_str,_ := ctx.JSONString(res)
	ctx.Write([]byte(json_str))
}

//代理 返利
func AgentRebate(U *userModel.User,money float64) error {
	//一级代理的返利
	err1 := U.CapitalUplete("+",FloatValue(money * 0.3,2),"用户返利")
	if err1 != nil {
		log.T("用户的一级代理返利失败！ error：%s",err1)
		return err1
	}
	//判断代理是否有上级
	level_two := userModel.GetUserById(uint32(U.ExtensionId))
	if level_two != nil {
		//找到代理的上级
		err2 := level_two.CapitalUplete("+",FloatValue(money * 0.07,2),"下级代理返利")
		if err2 != nil {
			log.T("代理的下一级代理返利失败！ error：%s",err2)
			return err2
		}
		level_three := userModel.GetUserById(uint32(level_two.ExtensionId))
		if level_three != nil {
			err3 := level_three.CapitalUplete("+",FloatValue(money * 0.07,2),"下级代理返利")
			if err3 != nil {
				log.T("代理的下一级代理返利失败！ error：%s",err3)
				return err3
			}
			return nil
		}
		return nil
	}
	return nil

}

func GetWeishu(str float64) int {
	var val int
	s :=fmt.Sprintf("%.2f", str)
	by := []byte(s)
	lengt := len(by)
	for i,_ := range by {
		if i == lengt-1{
			val,_ = strconv.Atoi(string(by[i]))
		}

	}
	return val
}

//返利给代理

func RebateLog(u *userModel.User,money float64) error {
	var err error = nil
	if u.ExtensionId != 0 {
		log.T("找到代理ID%d",u.ExtensionId)
		agent_info := userModel.GetUserById(uint32(u.ExtensionId))
		if agent_info != nil {
			log.T("获取到代理信息")
			err := agentModel.GetAgentRebateLog(agent_info.Id,u.Id,FloatValue(money*0.03,2))
			if err != nil {
				log.T("返利操作失败 %s",err)
			}
			return err
		}
		log.T("没有在数据库中找到代理信息")
		return err
	}
	log.T("没有代理")
	return err
}