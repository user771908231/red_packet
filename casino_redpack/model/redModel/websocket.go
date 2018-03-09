package redModel

import (
	"gopkg.in/olahol/melody.v1"
	"gopkg.in/mgo.v2/bson"
)
var WsServer = melody.New()

func init() {
	WsServer.HandleMessage(func(session *melody.Session, bytes []byte) {
		conn := WsConn(*session)
		conn.HandlMsg(bytes)
	})

	WsServer.HandleConnect(func(session *melody.Session) {
		conn := WsConn(*session)
		conn.HandConnect()
	})

}

//广播方法
func BroadCast(filter bson.M, dataFun func(conn WsConn)[]byte) {
	WsServer.BroadcastFilter([]byte{}, func(session *melody.Session) bool {
		conn := WsConn(*session)
		if len(filter) == 0 {
			session.Write(dataFun(conn))
			return true
		}
		for k, v := range filter {
			ex_v, ok := session.Get(k)
			if !ok || ex_v != v {
				return false
			}
		}
		session.Write(dataFun(conn))
		return true
	})
}

type WsConn melody.Session

//设置标签
func (conn WsConn) Set(key string, val interface{}) {
	session := conn.getWsSession()
	session.Set(key, val)
}

//获取标签
func (conn WsConn) Get(key string) (val interface{}) {
	session := conn.getWsSession()
	val,ex := session.Get(key)
	if ex == false {
		val = nil
	}
	return
}

//转换为melodySession
func (conn WsConn) getWsSession() melody.Session {
	var session melody.Session
	session = melody.Session(conn)
	return session
}

//广播给房间
func (conn WsConn) BroadToCurrRoom(msg []byte) {
	cur_roomId := conn.Get(TagRoomId)
	BroadCast(bson.M{TagRoomId: cur_roomId}, func(conn WsConn) []byte {
		return msg
	})
}

//给单个连接发消息
func (conn WsConn) WriteMsg(msg []byte) error {
	session := conn.getWsSession()
	return session.Write(msg)
}
//发送二进制数据
func (conn WsConn) WriteBinary(msg []byte) error {
	session := conn.getWsSession()
	return session.WriteBinary(msg)
}

//关闭连接
func (conn WsConn) Close() error {
	session := conn.getWsSession()
	return session.Close()
}
