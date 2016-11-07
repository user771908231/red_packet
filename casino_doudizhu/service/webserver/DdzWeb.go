package webserver

import (
	"net/http"
	"casino_server/common/log"
	"fmt"
)

func InitCms() {
	log.T("这里处理web请求")
	http.HandleFunc("/f", gameInfo) //设置访问的路由
	err := http.ListenAndServe(":9093", nil) //设置监听的端口 if err != nil {
	log.T("ListenAndServe: ", err)
}

func gameInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "\n打印朋友桌麻将的信息\n")
	fmt.Fprint(w, "\n朋友桌房间的信息打印完毕\n\n\n\n\n\n\n\n")
}
