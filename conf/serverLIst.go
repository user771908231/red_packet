package conf

import (
	"casino_common/common/consts/tableName"
	"casino_common/proto/ddproto"
	"casino_common/utils/db"
	"gopkg.in/mgo.v2/bson"
)

//服务器信息
type ServerInfo struct {
	GameId              ddproto.CommonEnumGame `bson:"GameId" title:"游戏Id"`
	Name                string `bson:"name" title:"游戏"`
	CurVersion          int32 `bson:"CurVersion" form:"CurVersion" binding:"" title:"当前版本"`
	IsUpdate            int32 `bson:"IsUpdate" form:"IsUpdate" binding:"" title:"是否需要强制升级"`
	IsMaintain          int32 `bson:"IsMaintain" form:"IsMaintain" binding:"" title:"是否在维护中"`
	MaintainMsg         string `bson:"MaintainMsg" form:"MaintainMsg" binding:"" title:"的维护信息"`
	ReleaseTag          int32 `bson:"ReleaseTag" form:"ReleaseTag" binding:"" title:"已经发布的版本"`
	DownloadUrl         string `bson:"DownloadUrl" form:"DownloadUrl" binding:"" title:"的下载连接"`
	LatestClientVersion int32 `bson:"LatestClientVersion" form:"LatestClientVersion" binding:"" title:"最后的版本号"`
	IP                  string `bson:"IP" form:"IP" binding:"" title:"ip"`
	PORT                int32 `bson:"PORT" form:"PORT" binding:"" title:"端口"`
	STATUS              int32 `bson:"STATUS" form:"STATUS" binding:"" title:"状态码"`
}

//游戏服务器的列表 都在这里
var ServerList []ServerInfo

func init() {
	err := db.C(tableName.DBT_GAME_CONFIG_LOGIN_LIST).FindAll(bson.M{}, &ServerList)
	if err != nil {
		panic("server list not found.")
		return
	}
}

//通过游戏id获取server info
func GetServerInfoByGameId(game_id ddproto.CommonEnumGame) *ServerInfo {
	for _, info := range ServerList {
		if info.GameId == game_id {
			return &info
		}
	}
	return nil
}

//获取登录服务器的ip地址
func GetAsLoginRpcAddress() string {
	return "192.168.2.188:44700"
}
