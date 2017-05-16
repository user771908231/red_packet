package config

import (
	"casino_admin/conf"
	"casino_admin/model/configModel"
	"casino_admin/modules"
	"casino_common/common/consts/tableName"
	"casino_common/common/log"
	"casino_common/common/rpc"
	"casino_common/common/rpc/protocol"
	"casino_common/common/service/configService"
	"casino_common/utils/db"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

//配置列表
func GameConfigListHandler(ctx *modules.Context) {
	table_name := ctx.Query("t")
	list, err := configService.GetConfig(table_name)
	if err != nil {
		ctx.Error("未找到该配置！", "", 0)
		return
	}
	field_list := configService.GetSliceField(list.List)
	type Id struct {
		Id bson.ObjectId `bson:"_id"`
	}
	id_list := []Id{}
	db.C(table_name).FindAll(bson.M{}, &id_list)

	ids_arr, _ := ctx.JSONString(id_list)
	ctx.Data["list"] = field_list
	ctx.Data["ids"] = ids_arr
	ctx.Data["table"] = table_name

	ctx.HTML(200, "admin/config/game/list")
}

type GameConfEditForm struct {
	Table string `binding:"Required"`
	Id    string `binding:"Required"`
	Field string `binding:"Required"`
	Value string `binding:"Required"`
}

//编辑字段
func GameConfigEdit(ctx *modules.Context) {
	id := ctx.Query("id")
	obj_id := bson.ObjectIdHex(id)

	result := configModel.GameConfigOne(obj_id)
	fmt.Println("success", obj_id)
	ctx.Data["config"] = result
	ctx.HTML(200, "admin/config/game/edit")
}

//编辑字段-更新
func GameConfigUpdate(ctx *modules.Context) {
	id := ctx.Query("id")
	obj_id := bson.ObjectIdHex(id)
	GameId := ctx.QueryFloat64("GameId")
	Name := ctx.Query("Name")
	CurVersion := ctx.QueryFloat64("CurVersion")
	IsUpdate := ctx.QueryFloat64("IsUpdate")
	IsMaintain := ctx.QueryFloat64("IsMaintain")
	MaintainMsg := ctx.Query("MaintainMsg")
	ReleaseTag := ctx.QueryFloat64("ReleaseTag")
	DownloadUrl := ctx.Query("DownloadUrl")
	LatestClientVersion := ctx.QueryFloat64("LatestClientVersion")
	IP := ctx.Query("IP")
	PORT := ctx.QueryFloat64("PORT")
	STATUS := ctx.QueryFloat64("STATUS")
	err := configModel.GameConfigUpdate(obj_id, GameId, Name, CurVersion, IsUpdate, IsMaintain, MaintainMsg, ReleaseTag, DownloadUrl, LatestClientVersion, IP, PORT, STATUS)
	if err == nil {
		rpc.Dial(conf.GetAsLoginRpcAddress(), rpc.AS_RELOAD_CONFIG, "", &protocol.CommonAckRpc{})
		ctx.Ajax(1, "编辑成功！", nil)
		//ctx.HTML(200,"admin/config/game/list")
	}
}

//编辑字段-更新
func GameConfigUpdateLogin(ctx *modules.Context) {
	log.T("开始更新登录服务器列表的信息...")
	id := ctx.Query("id")
	obj_id := bson.ObjectIdHex(id)
	CurVersion := ctx.QueryFloat64("CurVersion")
	BaseDownloadUrl := ctx.Query("BaseDownloadUrl")
	configModel.GameConfigUpdateLogin(obj_id, CurVersion, BaseDownloadUrl)
}

//登录服务器
func GameConfigList(ctx *modules.Context) {
	result := configModel.GameConfig()
	fmt.Println("success", result)
	ctx.Data["config"] = result
	ctx.HTML(200, "admin/config/game/list")
}

//登录服配置
func GameConfigLogin(ctx *modules.Context) {
	result := configModel.GameConfigLogin()
	fmt.Println("success", result)
	ctx.Data["config"] = result
	ctx.HTML(200, "admin/config/game/listLogin")
}

//新增一组配置
func GameConfigAddHandler(ctx *modules.Context) {
	table_name := ctx.Query("t")
	conf_info, err := configService.GetConfig(table_name)
	if err != nil {
		ctx.Error("未找到该配置！", "", 0)
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
	Data  string `binding:"Required"` //json表单
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
func GameServerInfoAddPost(ctx *modules.Context, form configService.LoginServerInfo) {
	err := db.C(tableName.DBT_GAME_CONFIG_LOGIN_LIST).Insert(form)
	if err != nil {
		ctx.Ajax(-1, "新增失败！", nil)
	}
	ctx.Ajax(1, "新增成功！", nil)
}
