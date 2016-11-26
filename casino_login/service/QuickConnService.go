package service

import "casino_login/conf"

var (
	GAME_ID_MJ int32 = 2; //斗地主
	GAME_ID_DDZ int32 = 3; //斗地主
)

//得到游戏的版本号
func GetLatestClientVersion(gameId int32) int32 {
	switch gameId {
	case GAME_ID_DDZ:
		return conf.Server.DDZ_LatestClientVersion
	case GAME_ID_MJ:
		return conf.Server.MJ_LatestClientVersion
	}
	return 0

}

func GetIsUpdate(gameId int32) int32 {
	switch gameId {
	case GAME_ID_DDZ:
		return conf.Server.DDZ_IsUpdate
	case GAME_ID_MJ:
		return conf.Server.MJ_IsUpdate
	}
	return 0
}

func GetDownloadUrl(gameId int32) string {
	switch gameId {
	case GAME_ID_DDZ:
		return conf.Server.DDZ_DownloadUrl
	case GAME_ID_MJ:
		return conf.Server.MJ_DownloadUrl
	}
	return conf.Server.BaseDownloadUrl
}

//是否需要维护
func GetIsMaintain(gameId int32) int32 {
	switch gameId {
	case GAME_ID_DDZ:
		return conf.Server.DDZ_IsMaintain
	case GAME_ID_MJ:
		return conf.Server.MJ_IsMaintain
	}
	return 0
}

func GetReleaseTagByVersion(gameId, v int32) int32 {
	var CurVersion int32 = 0

	//得到对应游戏的版本号
	switch gameId {
	case GAME_ID_DDZ:
		CurVersion = conf.Server.DDZ_ReleaseTag
	case GAME_ID_MJ:
		CurVersion = conf.Server.MJ_ReleaseTag
	}


	//返回游戏还是马甲
	if v <= CurVersion {
		return 1        //显示游戏
	} else {
		return 0        //显马甲
	}
}

//得到维护信息
func GetMaintainMsg(gameId int32) string {
	switch gameId {
	case GAME_ID_DDZ:
		return conf.Server.DDZ_MaintainMsg
	case GAME_ID_MJ:
		return conf.Server.MJ_MaintainMsg

	}
	return "服务器正在例行维护中!"
}

//游戏ip
func GetGameServerIp(gameId int32) string {
	switch gameId {
	case GAME_ID_DDZ:
		return conf.Server.DDZ_IP
	case GAME_ID_MJ:
		return conf.Server.MJ_IP
	}

	return "";
}

func GetGameServerPort(gameId int32) int32 {
	switch gameId {
	case GAME_ID_DDZ:
		return conf.Server.DDZ_PORT
	case GAME_ID_MJ:
		return conf.Server.MJ_PORT

	}

	return 8080;
}

func GetGameServerStatus(gameId int32) int32 {
	switch gameId {
	case GAME_ID_DDZ:
		return conf.Server.DDZ_STATUS
	case GAME_ID_MJ:
		return conf.Server.MJ_STATUS
	}
	return 1;
}
