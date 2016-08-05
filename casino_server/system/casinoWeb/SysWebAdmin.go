package casinoWeb

import (
	"net/http"
	"casino_server/common/log"
	"text/template"
	"casino_server/service/room"
)

func InitCms() {
	log.T("这里处理web请求")
	http.HandleFunc("/", gameInfo) //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口 if err != nil {
	log.T("ListenAndServe: ", err)
}


func gameInfo(w http.ResponseWriter,r *http.Request) {
	t,err := template.ParseFiles("system/casinoWeb/tmpl/gameinfo.html")
	if err != nil {
		log.E("errMsg:",err.Error())
	}
	p := room.ThGameRoomIns.ThDeskBuf[0]
	t.Execute(w,p)
}
