package Control

import (
	"casino_redpack/modules"
	"casino_redpack/model/userModel"
	"casino_redpack/model/Control"
	"gopkg.in/mgo.v2/bson"
)

func ShenzhizhongleiHandler(ctx *modules.Context) {
	uid := uint32(ctx.QueryInt("id"))
	msg := bson.M{
		"code":0,
		"msg":"设置失败！",
	}
	user_info := userModel.GetUserById(uid)
	if user_info != nil {
		C := Control.GetFindById(user_info.Id)
		if C != nil {
			C.Status.Open = 1
			err := C.Update()
			if err != nil {
				ctx.JSON(200,msg)
				return
			}
			msg["msg"] = "设置成功！"
			ctx.JSON(200,msg)
			return
		}
		C = new(Control.Control)
		C.Status.Open = 1
		err := C.Isert()
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

func ShenzhiredzhongleiHandler(ctx *modules.Context) {
	uid := uint32(ctx.QueryInt("id"))
	msg := bson.M{
		"code":1,
		"msg":"设置失败！",
	}
	user_info := userModel.GetUserById(uid)
	if user_info != nil {
		C := Control.GetFindById(user_info.Id)
		if C != nil {
			C.Status.Send = 1
			err := C.Update()
			if err != nil {
				ctx.JSON(200,msg)
				return
			}
			msg["msg"] = "设置成功！"
			ctx.JSON(200,msg)
			return
		}

		C = new(Control.Control)
		C.Status.Send = 1
		err := C.Isert()
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