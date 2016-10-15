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
		"麻将的信息：庄家的信息[%v]\t当前的游标[%v]：\n 麻将:\n %v"

		fmt.Fprintf(w, deskInfo, desk.GetDeskId(), desk.GetPassword(), desk.GetOwner(), desk.GetActiveUser(),
			desk.GetBanker(), desk.GetMJPaiCursor(), GetDeskMJInfo(desk))

		fmt.Fprintf(w, "\n开始打印user的信息:\n")
		for i, user := range desk.Users {
			if user != nil {
				fmt.Fprintf(w, "[%v],玩家的信息userId[%v],nickName[%v],status[%v],是否定缺[%v],定缺的花色[%v]\n", i, user.GetUserId(), "nickName", user.GetStatus(), user.GetDingQue(), user.GameData.HandPai.GetQueFlower())
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
		ii, _ := numUtils.Int2String(int32(i))
		pi, _ := numUtils.Int2String(p.GetValue())
		s = s + "\t (" + ii + "---" + pi + GetFlow(p.GetFlower()) + ")"
	}
	return s
}

func GetFlow(f int32) string {
	if f == 1 {
		return "铜"
	} else if f == 2 {
		return "条"
	} else if f == 3 {
		return "万"
	} else {
		return "白"
	}

}