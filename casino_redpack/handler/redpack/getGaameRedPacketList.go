package redpack

import (
	"casino_redpack/modules"
	"casino_redpack/model/userModel"
	"casino_redpack/model/redModel"
	"fmt"
)

//获取房间列表
func GetGaameRedPacketListHandler(ctx *modules.Context) {
	typeCode := ctx.Query("type")
	user := ctx.Session.Get("user")
	//	list := `{
	//	"code": 1,
	//	"message": "success",
	//	"request": []
	//}`
	var user_info userModel.User
	if user != nil {
		user_info = user.(userModel.User)
	}
	switch typeCode{
		case "1":

			//ctx.Write([]byte(list))

		case "2":
			fmt.Println("typeCode:"+typeCode)
		case "3":
			fmt.Println("typeCode:"+typeCode)
		case "4":
			fmt.Println("typeCode:"+typeCode)
		case "5":
			fmt.Println("typeCode:"+typeCode)
			redpacketLists := redModel.GetRedPacketRecord(user_info.Id)
			//Redpack := redModel.GetCreatorNameValues(user_info.Id)
			if redpacketLists != nil {
				list := redModel.GetLists(redpacketLists)
				ctx.Write([]byte(list))
			}
		default:
	}




//	list := `{
//	"code": 1,
//	"message": "success",
//	"request": [{
//		"id": 1137,
//		"type": 1,
//		"money": "50.00",
//		"all_membey": 5,
//		"has_member": 4,
//		"tail_number": 0,
//		"nickname": "zhujimeizu",
//		"headimgurl": "\/static\/userpic\/1190.jpg"
//	}, {
//		"id": 1138,
//		"type": 1,
//		"money": "100.00",
//		"all_membey": 5,
//		"has_member": 4,
//		"tail_number": 0,
//		"nickname": "\u4f60\u76f8\u8c8c\u5802\u5802",
//		"headimgurl": "\/static\/userpic\/536.jpg"
//	}]
//}`
//
//	ctx.Write([]byte(list))
}

func User_info (ctx *modules.Context) userModel.User{
	user := ctx.Session.Get("user")

	var user_info userModel.User
	user_info = user.(userModel.User)
	return  user_info

}
func  GetGaameRedPacketjlSendListHandler(ctx *modules.Context)  {
	typeCode := ctx.Query("type")
	switch typeCode {
	case "5":
		SendLists := redModel.GetPacketSendRecord(User_info(ctx).Id)
		if SendLists != nil {
			ctx.Write([]byte(SendLists))
		}else{
				list := `{
				"code": 0,
				"message": "success",
				"request": []
			}`
			ctx.Write([]byte(list))
		}

	}
}

func OpenPacketDetailsHandler(ctx *modules.Context)  {
	user := User_info(ctx)
	redId := ctx.QueryInt("redId")
	Id :=int32(redId)
	fmt.Println("ID",Id)
	value := redModel.OpenPacketDetails(Id,user.Id)
	ctx.Write([]byte(value))


}


