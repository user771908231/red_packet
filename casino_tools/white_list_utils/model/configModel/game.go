package configModel

import (
	"casino_common/utils/db"
	//"casino_common/common/consts"
	//"casino_tools/white_list_utils/modules"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/redisUtils"
)

const ADMIN_GAME_CONFIG string = "t_game_config_login_list"
const ADMIN_GAME_LOGIN string = "t_game_config_login"

type GameCongif struct {
	Id                  bson.ObjectId                         `bson:"_id"` //用户ID
	GameId              float64                        `bson:"GameId"`     //游戏ID
	Name                string                `bson:"name"`
	CurVersion          float64                `bson:"CurVersion"`
	IsUpdate            float64                `bson:"IsUpdate"`
	IsMaintain          float64                `bson:"IsMaintain"`
	MaintainMsg         string                `bson:"MaintainMsg"`
	ReleaseTag          float64                `bson:"ReleaseTag"`
	DownloadUrl         string                `bson:"DownloadUrl"`
	LatestClientVersion float64        `bson:"LatestClientVersion"`
	IP                  string                        `bson:"IP"` //进入IP
	PORT                float64                        `bson:"PORT"`
	STATUS              float64                        `bson:"STATUS"`
}
type GameCongifLogin struct {
	Id              bson.ObjectId                         `bson:"_id"` //用户ID
	CurVersion      float64                `bson:"CurVersion"`
	BaseDownloadUrl string                `bson:"BaseDownloadUrl"`
}

//登录服务器配置
func GameConfig() []*GameCongif {
	game_Congif := []*GameCongif{}
	err := db.C(ADMIN_GAME_CONFIG).FindAll(bson.M{}, &game_Congif)
	if err != nil {
		return nil
	}
	return game_Congif
}

//登录服配置
func GameConfigLogin() []*GameCongifLogin {
	game_Congif := []*GameCongifLogin{}
	err := db.C(ADMIN_GAME_LOGIN).FindAll(bson.M{}, &game_Congif)
	if err != nil {
		return nil
	}
	return game_Congif
}

//type GameList struct {
//	Id               *int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
//	gameID            []int32  `protobuf:"varint,2,opt,name=gameID" json:"gameID,omitempty"`
//	isOnMaintain          []int64  `protobuf:"varint,3,opt,name=isOnMaintain" json:"isOnMaintain,omitempty"`
//}

const PKEY_GAME_LIST string= "game_id"	//游戏服务器是否在维护中
const PKEY_USER_LIST string= "user_coin_redis_key"	//游戏服务器是否在维护中

//游戏服务器是否在维护中
func GameList() int64 {
	code := "14"
	return redisUtils.GetInt64(PKEY_GAME_LIST+ "_" +code)
}

func GameConfigOne(Id bson.ObjectId) *GameCongif {
	game_Congif := new(GameCongif)
	err := db.C(ADMIN_GAME_CONFIG).Find(bson.M{
		"_id": Id,
	}, game_Congif)
	if err != nil {
		return nil
	}
	return game_Congif
}

//登录服务器配置更新
func GameConfigUpdate(Id bson.ObjectId, GameId float64, Name string, CurVersion float64, IsUpdate float64, IsMaintain float64, MaintainMsg string, ReleaseTag float64, DownloadUrl string, LatestClientVersion float64, IP string, PORT float64, STATUS float64) (game_Congif *GameCongif) {
	err := db.C(ADMIN_GAME_CONFIG).Update(bson.M{"_id": Id}, bson.M{
		"$set": bson.M{
			"GameId":              GameId,
			"name":                Name,
			"CurVersion":          CurVersion,
			"IsUpdate":            IsUpdate,
			"IsMaintain":          IsMaintain,
			"MaintainMsg":         MaintainMsg,
			"ReleaseTag":          ReleaseTag,
			"DownloadUrl":         DownloadUrl,
			"LatestClientVersion": LatestClientVersion,
			"IP":                  IP,
			"PORT":                PORT,
			"STATUS":              STATUS,
		}, })
	if err != nil {
		return nil
	}
	return game_Congif
}

//登录服配置更新
func GameConfigUpdateLogin(Id bson.ObjectId, CurVersion float64, BaseDownloadUrl string) (game_Congif *GameCongif) {
	err := db.C(ADMIN_GAME_LOGIN).Update(bson.M{"_id": Id}, bson.M{
		"$set": bson.M{
			"CurVersion":      CurVersion,
			"BaseDownloadUrl": BaseDownloadUrl,
		}, })
	if err != nil {
		return nil
	}
	return game_Congif
}

//新增登录服
func GameConfigEditLogin(CurVersion float64, BaseDownloadUrl string) (game_Congif *GameCongif) {
	err := db.C(ADMIN_GAME_LOGIN).Insert(bson.M{
		"CurVersion" : CurVersion,
		"BaseDownloadUrl" : BaseDownloadUrl,
	})
	if err != nil {
		return nil
	}
	return game_Congif
}
