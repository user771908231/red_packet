package configModel

import (
	"gopkg.in/mgo.v2/bson"

	"casino_common/utils/db"
)
const ADMIN_GAME_CONFIG string = "t_game_config_login_list"
type GameCongif struct {
	Id          bson.ObjectId  			 `bson:"_id"`		//用户ID
	GameId      	float64 	 		`bson:"GameId"`		//游戏ID
	Name		string  		`bson:"name"`
	CurVersion	float64  		`bson:"CurVersion"`
	IsUpdate	float64  		`bson:"IsUpdate"`
	IsMaintain	float64  		`bson:"IsMaintain"`
	MaintainMsg	string  		`bson:"MaintainMsg"`
	ReleaseTag	float64  		`bson:"ReleaseTag"`
	DownloadUrl	string  		`bson:"DownloadUrl"`
	LatestClientVersion	float64  	`bson:"LatestClientVersion"`
	IP          string    			`bson:"IP"`			//进入IP
	PORT	float64  			`bson:"PORT"`
	STATUS	float64  			`bson:"STATUS"`
}
func GameConfigOne(Id bson.ObjectId) *GameCongif {
	game_Congif := new(GameCongif)
	err :=db.C(ADMIN_GAME_CONFIG).Find(bson.M{
		"_id" : Id,
	},game_Congif)
	if err != nil {
		return nil
	}
	return game_Congif
}

func GameConfigUpdate(Id bson.ObjectId,GameId float64,Name string,CurVersion float64,IsUpdate float64,IsMaintain float64,MaintainMsg string,ReleaseTag float64,DownloadUrl string,LatestClientVersion float64,IP string,PORT float64,STATUS float64) (game_Congif *GameCongif) {
	err := db.C(ADMIN_GAME_CONFIG).Update(bson.M{"_id": Id},bson.M{
		"$set" : bson.M{
			"GameId" : GameId,
			"name" : Name,
			"CurVersion" : CurVersion,
			"IsUpdate" : IsUpdate,
			"IsMaintain" : IsMaintain,
			"MaintainMsg" : MaintainMsg,
			"ReleaseTag" : ReleaseTag,
			"DownloadUrl" : DownloadUrl,
			"LatestClientVersion" : LatestClientVersion,
			"IP" : IP,
			"PORT" : PORT,
			"STATUS" : STATUS,
		},})
	if err != nil {
		return nil
	}
	return game_Congif
}
