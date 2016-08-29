package casinoWeb

import (
	"net/http"
	"casino_server/common/log"
)

func InitCms() {
	log.T("这里处理web请求")
	http.HandleFunc("/", gameInfo) //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口 if err != nil {
	log.T("ListenAndServe: ", err)
}


func gameInfo(w http.ResponseWriter,r *http.Request) {
	//w.Write([]byte("朋友桌房间的信息:\n"))
	//for i:=0;i<room.ThGameRoomIns ;i++ {
	//	desk := room.ThGameRoomIns.ThDeskBuf[i]
	//	deskTitle := "开始打印desk.id[%v],roomKey[%v]的信息:\n"+"sfsfl"
	//	fmt.Printf(deskTitle,desk.Id,desk.RoomKey)
	//
	//	//w.Write([]byte(deskTitle)
	//
	//}
}