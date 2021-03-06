package config

import (
	"casino_redpack/model/configModel"
	"casino_redpack/modules"
	"casino_common/common/consts/tableName"
	"casino_common/common/log"
	"casino_common/common/service/configService"
	"casino_common/utils/db"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"casino_common/utils/redisUtils"
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
		//rpc.Dial(conf.GetAsLoginRpcAddress(), rpc.AS_RELOAD_CONFIG, "", &protocol.CommonAckRpc{})
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
	ctx.Data["config"] = result
	ctx.HTML(200, "admin/config/game/list")
}


//游戏服务器
func GameListHandler(ctx *modules.Context){
	//code := ctx.Query("code")
	ctx.Data["info_1"] = redisUtils.GetInt64("game_id_1")
	ctx.Data["info_2"] = redisUtils.GetInt64("game_id_2")
	ctx.Data["info_3"] = redisUtils.GetInt64("game_id_3")
	ctx.Data["info_4"] = redisUtils.GetInt64("game_id_4")
	ctx.Data["info_5"] = redisUtils.GetInt64("game_id_5")
	ctx.Data["info_7"] = redisUtils.GetInt64("game_id_7")
	ctx.Data["info_8"] = redisUtils.GetInt64("game_id_8")
	ctx.Data["info_9"] = redisUtils.GetInt64("game_id_9")
	ctx.Data["info_10"] = redisUtils.GetInt64("game_id_10")
	ctx.Data["info_11"] = redisUtils.GetInt64("game_id_11")
	ctx.Data["info_12"] = redisUtils.GetInt64("game_id_12")
	ctx.Data["info_13"] = redisUtils.GetInt64("game_id_13")
	ctx.Data["info_14"] = redisUtils.GetInt64("game_id_14")
	ctx.Data["info_15"] = redisUtils.GetInt64("game_id_15")
	ctx.Data["info_16"] = redisUtils.GetInt64("game_id_16")
	ctx.Data["info_17"] = redisUtils.GetInt64("game_id_17")
	ctx.Data["info_18"] = redisUtils.GetInt64("game_id_18")
	ctx.Data["info_19"] = redisUtils.GetInt64("game_id_19")
	ctx.Data["info_20"] = redisUtils.GetInt64("game_id_20")
	ctx.HTML(200, "admin/config/game/game_server")
}

func GameListEditHandler(ctx *modules.Context) {
	code := ctx.Query("code")
	edit := ctx.QueryInt64("edit")
	redisUtils.SetInt64("game_id_"+code,edit)
	ctx.Success("操作成功！", "/admin/config/game/gameList", 1)
}


//登录服配置
func GameConfigLogin(ctx *modules.Context) {
	result := configModel.GameConfigLogin()
	fmt.Println("success", result)
	ctx.Data["config"] = result
	ctx.HTML(200, "admin/config/game/listLogin")
}

//新增登录服
func GameConfigEditLogin(ctx *modules.Context) {
	CurVersion := ctx.QueryFloat64("CurVersion")
	BaseDownloadUrl := ctx.Query("BaseDownloadUrl")
	configModel.GameConfigEditLogin(CurVersion,BaseDownloadUrl)

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
