package webserver

import (
	"net/http"
	"casino_server/common/log"
	"fmt"
	"casino_majiang/service/majiang"
	"casino_server/utils/numUtils"
)

func InitCms() {
	log.T("这里处理web请求")
	http.HandleFunc("/f", gameInfo) //设置访问的路由
	err := http.ListenAndServe(":9091", nil) //设置监听的端口 if err != nil {
	log.T("ListenAndServe: ", err)
}

func gameInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "\n打印朋友桌麻将的信息\n")
	for _, desk := range majiang.FMJRoomIns.Desks {
		printDeskInfo(w, desk)
	}
	fmt.Fprint(w, "\n朋友桌房间的信息打印完毕\n\n\n\n\n\n\n\n")
}

func printDeskInfo(w http.ResponseWriter, desk *majiang.MjDesk) {
	if desk != nil {
		deskInfo := "开始打印desk.id[%v], \t房间号roomKey[%v]的信息:\t 房主Owner[%v],activeUser[%v]\t \n" +
		"麻将的信息：庄家的信息[%v]\t当前的游标[%v]：\n 麻将:\n %v checkCase:%v \n"

		fmt.Fprintf(w, deskInfo, desk.GetDeskId(), desk.GetPassword(), desk.GetOwner(), desk.GetActiveUser(),
			desk.GetBanker(), desk.GetMJPaiCursor(), GetDeskMJInfo(desk), desk.CheckCase)

		fmt.Fprintf(w, "\n开始打印user的信息:\n")
		for i, user := range desk.Users {
			if user != nil {
				fmt.Fprintf(w, "[%v],玩家的信息userId[%v],nickName[%v],status[%v],是否定缺[%v],定缺的花色[%v]\n 玩家的手牌[%v]\n玩家的碰牌[%v],玩家的杠牌[%v],玩家的胡牌[%v],玩家的inpai[%v]\n",
					i, user.GetUserId(), "nickName", user.GetStatus(), user.GetDingQue(), user.GameData.HandPai.GetQueFlower(), getUserPaiInfo(user),
					getUserPengPaiInfo(user),
					getUserGnagPaiInfo(user),
					getUserHuPaiInfo(user), getUserInPaiInfo(user))
			}
		}

		fmt.Fprintf(w, "打印desk.id[%v]完毕 \n\n\n\n\n\n", desk.GetDeskId())
	}

}

func GetDeskMJInfo(desk *majiang.MjDesk) string {
	if desk == nil || desk.AllMJPai == nil {
		return "暂时没有初始化麻将"
	}
	s := ""
	for i, p := range desk.AllMJPai {
		is, _ := numUtils.Int2String(int32(i))

		ii, _ := numUtils.Int2String(int32(p.GetIndex()))
		s = s + " (" + is + "-" + ii + "-" + p.LogDes() + ")"
	}
	return s
}

func getUserPaiInfo(user *majiang.MjUser) string {
	if user.GameData == nil || user.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	s := ""
	for _, p := range user.GameData.HandPai.Pais {
		ii, _ := numUtils.Int2String(int32(p.GetIndex()))
		s = s + ii + "-" + p.LogDes() + "\t "
	}

	return s

}

func getUserPengPaiInfo(user *majiang.MjUser) string {
	if user.GameData == nil || user.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	s := ""
	for _, p := range user.GameData.HandPai.PengPais {
		ii, _ := numUtils.Int2String(int32(p.GetIndex()))
		s = s + ii + "-" + p.LogDes() + "\t "
	}

	return s

}

func getUserGnagPaiInfo(user *majiang.MjUser) string {
	if user.GameData == nil || user.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	s := ""
	for _, p := range user.GameData.HandPai.GangPais {
		ii, _ := numUtils.Int2String(int32(p.GetIndex()))
		s = s + ii + "-" + p.LogDes() + "\t "
	}

	return s

}

func getUserHuPaiInfo(user *majiang.MjUser) string {
	if user.GameData == nil || user.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	s := ""
	for _, p := range user.GameData.HandPai.HuPais {
		ii, _ := numUtils.Int2String(int32(p.GetIndex()))
		s = s + ii + "-" + p.LogDes() + "\t "
	}

	return s

}
func getUserInPaiInfo(user *majiang.MjUser) string {
	if user.GameData == nil || user.GameData.HandPai == nil {
		return "用户还没有牌"
	}

	return user.GameData.HandPai.InPai.LogDes()

}
