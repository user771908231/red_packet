package casinoWeb

import (
	"net/http"
	"casino_server/common/log"
	"casino_server/service/room"
	"fmt"
)

func InitCms() {
	log.T("这里处理web请求")
	http.HandleFunc("/f", gameInfo) //设置访问的路由
	http.HandleFunc("/c", gameInfocs) //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口 if err != nil {
	log.T("ListenAndServe: ", err)
}

func gameInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "朋友桌房间的信息:\n")
	for i := 0; i < len(room.ThGameRoomIns.ThDeskBuf); i++ {
		desk := room.ThGameRoomIns.ThDeskBuf[i]
		if desk != nil {
			deskInfo := "开始打印desk.id[%v],roomKey[%v]的信息:\n" +
			"房主:%v\n" +
			"游戏类型:%v\n" +
			"庄家:%v\n" +
			"小盲注:%v\n" +
			"大盲注:%v\n" +
			"带入金额:%v\n" +
			"第几轮:%v\n" +
			"第几局:%v\n" +
			"总局数:%v\n" +
			"主奖池:%v\n" +
			"边池:%v\n" +
			"最低加注额度:%v\n" +
			"总游戏玩家:%v\n" +
			"已经准备的玩家:%v\n" +
			"------------------所有玩家的信息:\n"
			fmt.Fprintf(w, deskInfo, desk.Id, desk.RoomKey, desk.DeskOwner, desk.GameType, desk.SmallBlind, desk.BigBlind,
				desk.InitRoomCoin, desk.RoundCount, desk.JuCountNow, desk.JuCount, desk.Jackpot, desk.EdgeJackpot, desk.MinRaise,
				desk.UserCount, desk.GetGameReadyCount())

			for j := 0; j < len(desk.Users); j++ {
				u := desk.Users[j]
				if u != nil {
					userInfo := "当前desk[%v]的user[%v],seatId[%v],nickname[%v]的状态status[%v],HandCoin[%v],TurnCoin[%v],RoomCoin[%v],isBreak[%v],isLeave[%v]\n"
					fmt.Fprintf(w, userInfo, desk.Id, u.UserId, u.Seat, u.NickName, u.Status, u.HandCoin, u.TurnCoin, u.RoomCoin, u.IsBreak, u.IsLeave)
				}
			}
			fmt.Fprint(w, "------------------所有玩家的信息打印完毕\n", desk.Id)

			fmt.Fprint(w, "desk[%v]的信息打印完毕:\n\n\n\n", desk.Id)

		}
	}
	fmt.Fprint(w, "朋友桌房间的信息打印完毕:\n")
}

func gameInfocs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "锦标赛房间的信息:\n")
	for i := 0; i < len(room.ChampionshipRoom.ThDeskBuf); i++ {
		desk := room.ThGameRoomIns.ThDeskBuf[i]
		if desk != nil {
			deskInfo := "开始打印desk.id[%v],roomKey[%v]的信息:\n" +
			"房主:%v\n" +
			"游戏类型:%v\n" +
			"庄家:%v\n" +
			"小盲注:%v\n" +
			"大盲注:%v\n" +
			"带入金额:%v\n" +
			"第几轮:%v\n" +
			"第几局:%v\n" +
			"总局数:%v\n" +
			"主奖池:%v\n" +
			"边池:%v\n" +
			"最低加注额度:%v\n" +
			"总游戏玩家:%v\n" +
			"已经准备的玩家:%v\n" +
			"------------------所有玩家的信息:\n"
			fmt.Fprintf(w, deskInfo, desk.Id, desk.RoomKey, desk.DeskOwner, desk.GameType, desk.SmallBlind, desk.BigBlind,
				desk.InitRoomCoin, desk.RoundCount, desk.JuCountNow, desk.JuCount, desk.Jackpot, desk.EdgeJackpot, desk.MinRaise,
				desk.UserCount, desk.GetGameReadyCount())

			for j := 0; j < len(desk.Users); j++ {
				u := desk.Users[j]
				if u != nil {
					userInfo := "当前desk[%v]的user[%v],seatId[%v],nickname[%v]的状态status[%v],HandCoin[%v],TurnCoin[%v],RoomCoin[%v],isBreak[%v],isLeave[%v]\n"
					fmt.Fprintf(w, userInfo, desk.Id, u.UserId, u.Seat, u.NickName, u.Status, u.HandCoin, u.TurnCoin, u.RoomCoin, u.IsBreak, u.IsLeave)
				}
			}
			fmt.Fprint(w, "------------------所有玩家的信息打印完毕\n", desk.Id)

			fmt.Fprint(w, "desk[%v]的信息打印完毕:\n", desk.Id)

		}
	}
	fmt.Fprint(w, "锦标赛房间的信息打印完毕:\n\n\n\n")
}