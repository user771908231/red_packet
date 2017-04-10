package config

import (
	"casino_super/modules"
	"casino_common/common/service/configService"
	"casino_common/utils/db"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"encoding/json"
	login_conf "casino_login/conf"
	"casino_common/common/consts/tableName"
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

	ids_arr,_ := ctx.JSONString(id_list)
	ctx.Data["list"] = field_list
	ctx.Data["ids"] = ids_arr
	ctx.Data["table"] = table_name

	ctx.HTML(200, "admin/config/game/list")
}

type GameConfEditForm struct {
	Table string `binding:"Required"`
	Id string `binding:"Required"`
	Field string `binding:"Required"`
	Value string `binding:"Required"`
}
//编辑字段
func GameConfigEditPost(ctx *modules.Context, form GameConfEditForm) {
	err := db.C(form.Table).Update(bson.M{
		"_id": bson.ObjectIdHex(form.Id),
	}, bson.M{
		"$set": bson.M{form.Field: form.Value},
	})
	if err != nil {
		ctx.Ajax(-1, "更新失败！Error:"+err.Error(), err.Error())
		return
	}
	ctx.Ajax(1, "更新成功！", nil)
}

//新增一组配置
func GameConfigAddHandler(ctx *modules.Context) {
	table_name := ctx.Query("t")
	conf_info, err := configService.GetConfig(table_name)
	if err != nil {
		ctx.Error("未找到该配置！","", 0)
		return
	}

	list := [][]configService.FieldInfo{}
	list = append(list, configService.GetColInfo(conf_info.Row))
	ctx.Data["table"] = table_name
	ctx.Data["list"] = list
	ctx.HTML(200, "admin/config/game/add")
}

//新增一组配置
type GameConfigAddForm struct {
	Table string `binding:"Required"`
	Data string `binding:"Required"` //json表单
}

func GameConfigAddPost(ctx *modules.Context, form GameConfigAddForm) {
	conf_info, err := configService.GetConfig(form.Table)
	if err != nil {
		ctx.Ajax(-2, "未找到该表单！", nil)
		return
	}
	new_row := reflect.New(reflect.TypeOf(conf_info.Row).Elem()).Interface()

	err = json.Unmarshal([]byte(form.Data), new_row)

	if err != nil {
		ctx.Ajax(-1, "新增失败！Error:"+err.Error(), err.Error())
		return
	}

	ctx.Ajax(1, "", new_row)
}


//新增ServerInfo
func GameServerInfoAddPost(ctx *modules.Context, form login_conf.ServerInfo) {
	err := db.C(tableName.DBT_GAME_CONFIG_LOGIN_LIST).Insert(form)
	if err != nil {
		ctx.Ajax(-1, "新增失败！", nil)
	}
	ctx.Ajax(1, "新增成功！", nil)
}
