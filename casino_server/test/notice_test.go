package mongodb

import (
	"net"
	"testing"
	"github.com/name5566/leaf/db/mongodb"
	"fmt"
	"casino_server/mode"
	"majiang/conf/db"
	"casino_server/service/noticeServer"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/test"
)

func TestNotice(t *testing.T) {
	initSys()
	//saveNotice()
	getNoticeByType(noticeServer.NOTICE_TYPE_GUNDONG)		//type
	//getNoticeByType(noticeServer.NOTICE_TYPE_CHONGZHI)		//type
	//getNoticeByType(noticeServer.NOTICE_TYPE_GONGGAO)		//type
}

func saveNotice(){

	mongc, err := mongodb.Dial(db.DB_IP, db.DB_PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer mongc.Close()

	// 获取回话 session
	s := mongc.Ref()
	defer mongc.UnRef(s)

	saveNotice :=  &mode.T_th_notice{}

	saveNotice.Id = 1
	saveNotice.NoticeType = noticeServer.NOTICE_TYPE_CHONGZHI
	saveNotice.NoticeContent = "测试公告"
	saveNotice.NoticeTitle	 = "测试公告的标题"
	saveNotice.NoticeFileds  = []string{"微信","wx001","微信","wx002","qq","qq001"}
	noticeServer.SaveNotice(saveNotice)
}

//通过类型来得到公告
func getNoticeByType(noticeType int32){
	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	pid := int32(bbproto.EProtoId_PID_GAME_GAMENOTICE)
	data2 := &bbproto.GameNotice{}
	data2.NoticeType = &noticeType
	m2 := test.AssembleDataNomd5(uint16(pid), data2)
	conn.Write(m2)
	result := test.Read(conn).(*bbproto.Game_AckNotice)
	fmt.Println("notice.title:",result.GetNoticeTitle())
	fmt.Println("notice.fs:",result.GetFileds())

}


