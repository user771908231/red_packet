package config

import (
	"casino_super/modules"
	"casino_common/common/service/configService"
	"casino_common/utils/db"
	"gopkg.in/mgo.v2/bson"
)

//配置列表
func GameConfigListHandler(ctx *modules.Context) {
	table_name := ctx.Query("t")
	list, err := configService.GetConfig(table_name)
	if err != nil {
		ctx.Error("未找到该配置！","", 0)
		return
	}
	field_list := configService.GetSliceField(list.List)
	type Id struct {
		Id bson.ObjectId `bson:"_id"`
	}
	id_list := []Id{}
	db.C(table_name).FindAll(bson.M{}, &id_list)

	ctx.Data["list"] = field_list
	ctx.Data["ids"] = id_list

	ctx.HTML(200, "admin/config/game/list")
}
