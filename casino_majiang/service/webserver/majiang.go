package webserver

import (
	"net/http"
	"casino_server/common/log"
	"fmt"
	"casino_majiang/service/majiang"
)

func InitCms() {
	log.T("这里处理web请求")
	http.HandleFunc("/f", gameInfo) //设置访问的路由
	err := http.ListenAndServe(":9091", nil) //设置监听的端口 if err != nil {
	log.T("ListenAndServe: ", err)
}

func gameInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "麻将....朋友桌房间的信息:\n")

	for _, desk := range majiang.FMJRoomIns.Desks {
		printDeskInfo(w, desk)
	}
	fmt.Fprint(w, "朋友桌房间的信息打印完毕:\n")
}

func printDeskInfo(w http.ResponseWriter, desk *majiang.MjDesk) {
	if desk != nil {
		deskInfo := "开始打印desk.id[%v],roomKey[%v]的信息:\n" +
		"------------------打印Users的信息:\n"
		fmt.Fprintf(w, deskInfo, desk.GetDeskId(), desk.GetPassword())
		fmt.Fprint(w, "打印完毕:\n", desk.GetDeskGameInfo())
	}
}


