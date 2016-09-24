package casinoWeb

import (
	"net/http"
	"casino_server/common/log"
	"casino_server/service/room"
	"fmt"
	"casino_server/utils/timeUtils"
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
		printDeskInfo(w, desk)
	}
	fmt.Fprint(w, "朋友桌房间的信息打印完毕:\n")
}

func gameInfocs(w http.ResponseWriter, r *http.Request) {
	fr := room.GetFirstCSTHGame()
	fmt.Fprint(w, "锦标赛房间matchId[%v],盲注等级[%v],readyTime[%v],beginTime[%v],endTime[%v],\n rankinfo[%v]",
		fr.MatchId, fr.BlindLevel, timeUtils.Format(fr.ReadyTime), timeUtils.Format(fr.BeginTime), timeUtils.Format(fr.EndTime), fr.RankInfo)
	for i := 0; i < len(fr.ThDeskBuf); i++ {
		desk := fr.ThDeskBuf[i]
		printDeskInfo(w, desk)
	}
	fmt.Fprint(w, "\n锦标赛房间的信息打印完毕:\n\n\n\n")
}

func printDeskInfo(w http.ResponseWriter, desk *room.ThDesk) {
	if desk != nil {
		deskInfo := "开始打印desk.id[%v],roomKey[%v]的信息:\n" +
		"房主:%v  游戏类型:%v 带入金额:%v  总局数:%v\n " +
		"庄家:%v  小盲注:%v  大盲注:%v\n" +
		"第几轮:%v  第几局:%v\n" +
		"sendFlop[%v]	,sendTurn[%v],	sendRive[%v]" +
		"主奖池:%v  边池:%v 最低加注额度:%v\n" +
		"总游戏玩家:%v 已经准备的玩家:%v getlotteryCheckFalseCount:%v \n" +
		"RaiseUserId :%v  当前轮到谁出牌:BetUserNow[%v]\n" +
		"------------------打印Users的信息:\n"
		fmt.Fprintf(w, deskInfo,
			desk.Id, desk.RoomKey,
			desk.DeskOwner, desk.GameType, desk.InitRoomCoin, desk.JuCount,
			desk.Dealer, desk.SmallBlind, desk.BigBlind,
			desk.RoundCount, desk.JuCountNow,
			desk.SendFlop, desk.SendTurn, desk.SendRive,
			desk.Jackpot, desk.EdgeJackpot, desk.MinRaise,
			desk.GetUserCount(), desk.GetGameReadyCount(), desk.GetLotteryCheckFalseCount(),
			desk.RaiseUserId, desk.BetUserNow)

		for j := 0; j < len(desk.Users); j++ {
			u := desk.Users[j]
			printThUserInfo(w, desk, u)        //打印user的信息
		}

		fmt.Fprint(w, "------------------打印leaveUsers的信息\n", desk.Id)
		for j := 0; j < len(desk.LeaveUsers); j++ {
			u := desk.LeaveUsers[j]
			printThUserInfo(w, desk, u)        //打印user的信息
		}

		fmt.Fprint(w, "------------------所有玩家的信息打印完毕\n", desk.Id)

		fmt.Fprint(w, "desk[%v]的信息打印完毕:\n\n\n\n", desk.Id)

	}
}

func printThUserInfo(w http.ResponseWriter, desk *room.ThDesk, u *room.ThUser) {
	if u != nil {
		userInfo := "当前desk[%v]的user[%v],seatId[%v],nickname[%v]的状态status[%v],HandCoin[%v],TurnCoin[%v],RoomCoin[%v],isBreak[%v],isLeave[%v],LotteryCheck[%v],u.gameStatus[%v],csgaminStatus[%v]-\n"
		fmt.Fprintf(w, userInfo, desk.Id, u.UserId, u.Seat, u.NickName, u.GetStatusDes(), u.HandCoin, u.TurnCoin, u.RoomCoin, u.IsBreak, u.IsLeave, u.LotteryCheck, u.GameStatus, u.CSGamingStatus)

	}

}

