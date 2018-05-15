package Control

import (
	"casino_redpack/modules"
	"casino_redpack/model/userModel"
	"casino_redpack/model/Control"
	"gopkg.in/mgo.v2/bson"

	"casino_common/common/log"
)

func ShenzhizhongleiHandler(ctx *modules.Context) {
	uid := uint32(ctx.QueryInt("id"))
	log.T("----------------------------------------------uid:",uid)
	msg := bson.M{
		"code":0,
		"msg":"设置失败！",
	}
	user_info := userModel.GetUserById(uid)
	log.T("user",user_info)
	if user_info != nil {
		C := Control.GetFindById(user_info.Id)
		if C != nil {
			log.T("在控制表中找到--开始更新")
			C.Status.Open = 1
			err := C.Update()
			if err != nil {
				log.T("更新一条失败")
				ctx.JSON(200,msg)
				return
			}
			msg["msg"] = "设置成功！"
			ctx.JSON(200,msg)
			return
		}
		log.T("没有在控制表中找到--开始新家一条")
		C1 := new(Control.Control)
		C1.UserId = user_info.Id
		C1.Status.Open = 1
		err := C1.Isert()
		if err != nil {
			log.T("新家一条失败")
			ctx.JSON(200,msg)
			return
		}
		msg["msg"] = "设置成功！"
		ctx.JSON(200,msg)
		return
	}
	ctx.JSON(200,msg)
}

func ShenzhiredzhongleiHandler(ctx *modules.Context) {
	uid := uint32(ctx.QueryInt("id"))
	log.T("----------------------------------------------uid:",uid)
	msg := bson.M{
		"code":0,
		"msg":"设置失败！",
	}
	user_info := userModel.GetUserById(uid)
	if user_info != nil {
		C := Control.GetFindById(user_info.Id)
		if C != nil {
			C.Status.Send = 0
			err := C.Update()
			if err != nil {
				ctx.JSON(200,msg)
				return
			}
			msg["msg"] = "设置成功！"
			ctx.JSON(200,msg)
			return
		}

		C1 := new(Control.Control)
		C1.Status.Send = 1
		C1.UserId = user_info.Id
		err := C1.Isert()
		if err != nil {
			ctx.JSON(200,msg)
			return
		}
		msg["msg"] = "设置成功！"
		ctx.JSON(200,msg)
		return
	}
	ctx.JSON(200,msg)
}

func ShenzhengchangHandler(ctx *modules.Context) {
	uid := uint32(ctx.QueryInt("id"))
	msg := bson.M{
		"code":1,
		"msg":"设置失败！",
	}
	user_info := userModel.GetUserById(uid)
	if user_info != nil {
		C := Control.GetFindById(user_info.Id)
		if C != nil {
			err :=C.Del()
			if err != nil {
				ctx.JSON(200,msg)
				return
			}
		}
		msg["msg"] = "设置成功！"
		ctx.JSON(200,msg)
		return
	}
	ctx.JSON(200,msg)
}