package deskHandler

import (
	"casino_clientlog/modules"
	"time"
	"casino_common/utils/numUtils"
	"casino_common/proto/ddproto"
	"casino_common/common/consts/tableName"
	"casino_common/utils/db"
	"gopkg.in/mgo.v2/bson"
)

//输赢
type DeskWinItem struct {
	UserId    uint32
	NickName  string
	WinAmount int64
}
//单条数据
type DeskUsersRow struct {
	DeskId     int32
	Password   string
	GameNumber int32
	UserIds    string
	BeginTime  time.Time
	EndTime    time.Time
	Records    []DeskWinItem
	RoundStr   string //局数信息
}

//表单
func FindUsersHandler(ctx *modules.Context) {
	passwd := ctx.Query("passwd")
	gid := ddproto.CommonEnumGame_GID_SRC
	table := ""
	list := []DeskUsersRow{}

	if len(passwd) == 6 {
		gid_str := string(passwd[2]) + string(passwd[3])
		gid_int := numUtils.String2Int(gid_str)
		gid = ddproto.CommonEnumGame(int32(gid_int))

		switch gid {
		case ddproto.CommonEnumGame_GID_DDZ:
			table = tableName.DBT_DDZ_DESK_ROUND_ALL
		case ddproto.CommonEnumGame_GID_MAHJONG:
			table = tableName.DBT_MJ_DESK_ROUND_ALL
		case ddproto.CommonEnumGame_GID_ZXZ:
			table = tableName.DBT_MJ_ZXZ_DESK_ROUND_ALL
		case ddproto.CommonEnumGame_GID_PDK:
			table = tableName.DBT_PDK_DESK_ROUND_ALL
		case ddproto.CommonEnumGame_GID_NIUNIUJINGDIAN:
			table = tableName.DBT_NIU_DESK_ROUND_ALL
		case ddproto.CommonEnumGame_GID_BANTUOZI:
			table = tableName.DBT_PEZ_DESK_ROUND_ALL
		case ddproto.CommonEnumGame_GID_MJBAISHAN:
			table = tableName.DBT_MJ_BS_DESK_ROUND_ALL
		}

		db.Log(table).FindAll(bson.M{
			"password": passwd,
		}, &list)

	}

	//ctx.Dump(list)
	ctx.Data["list"] = list
	ctx.Data["passwd"] = passwd
	ctx.HTML(200, "log/desk_users")

}
