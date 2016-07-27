package testService

import (
	"casino_server/msg/bbprotogo"
	"casino_server/common/log"
	"casino_server/service/room"
)

func Handlertestp(m *bbproto.TestP1){
	ttype := m.GetName2()
	log.T("-----------测试type[%v]--------",ttype)
	switch ttype {
	case "ogg":
		//答应当前房间的信息
	printGameinfo()
	case "":
		log.T("什么都没有")

	}
}

func printGameinfo(){
	game := room.ThGameRoomIns
	//循环desk打印信息
	log.T("-----------打印游戏信息--------")


	for i := 0; i < len(game.ThDeskBuf); i++ {
		d := game.ThDeskBuf[i]
		log.T("-----------打印desk[%v]的信息,begin--------" ,i)
		if d != nil {
			log.T("-----------desk[%v]的状态status[%v]--------" ,i,d.Status)
			log.T("-----------desk[%v]的庄 Dealer[%v]--------" ,i,d.Dealer)
			log.T("-----------desk[%v]的押注BetUserNow[%v]--------" ,i,d.BetUserNow)

			log.T("-----------打印desk[%v]的用户的信息,begin--------" ,i)
			for j := 0; j < len(d.Users); j++ {
				u := d.Users[j]
				log.T("-----------打印desk[%v]的用户[%v]的信息,begin--------" ,i,j)

				if u == nil {
					log.T("-----------打印desk[%v]的用户[%v]==nil--------" ,i,j)

				}else{
					log.T("-----------打印desk[%v]的用户[%v]的id[%v]--------" ,i,j,u.UserId)
					log.T("-----------打印desk[%v]的用户[%v]的HandCoin[%v]--------" ,i,j,u.HandCoin)
					log.T("-----------打印desk[%v]的用户[%v]Coin[%v]--------" ,i,j,u.Coin)
					log.T("-----------打印desk[%v]的用户[%v]cards[%v]--------" ,i,j,u.Cards)
					log.T("-----------打印desk[%v]的用户[%v]iStatus[%v]--------" ,i,j,u.Status)
					log.T("-----------打印desk[%v]的用户[%v]TotalBet[%v]--------" ,i,j,u.TotalBet)
					log.T("-----------打印desk[%v]的用户[%v]TurnCoin[%v]--------" ,i,j,u.TurnCoin)
					log.T("-----------打印desk[%v]的用户[%v]NickName[%v]--------" ,i,j,u.NickName)
				}

				log.T("-----------打印desk[%v]的用户[%v]的信息,end----------" ,i,j)
			}
			log.T("-----------打印desk[%v]的用户的信息,end----------" ,i)


		}else{
			log.T("-----------desk[%v]==nil--------" ,i)
		}
		log.T("-----------打印desk[%v]的信息,end----------" ,i)
	}

	log.T("-----------打印游戏信息结束--------")
}
