package test

import (
	"testing"
	"casino_server/common/log"
	"casino_majiang/msg/funcsInit"
	"github.com/golang/protobuf/proto"
	"net"
	"encoding/binary"
)

func TestLoginIn(t *testing.T) {
	log.Debug("login in testing")
	//创建一个连接
	conn, err := net.Dial("tcp", "192.168.199.152:3799")
	if err != nil {
		panic(err)
	}

	gameAckQuickConn := newProto.NewGame_AckQuickConn()
	*gameAckQuickConn.CurrVersion = 1
	*gameAckQuickConn.DownloadUrl = ""
	//gameAckQuickConn.GameServer =
	gameAckQuickConn.Header = newProto.NewHeader()
	*gameAckQuickConn.IsMaintain = 1
	//*gameAckQuickConn.MaintainMsg = ""
	*gameAckQuickConn.IsUpdate = 1
	*gameAckQuickConn.ReleaseTag = 1
	*gameAckQuickConn.VersionInfo = ""

	//proto 编码
	data, _ := proto.Marshal(gameAckQuickConn)
	log.D("data is :%v", data)
	// len + data
	m := make([]byte, 2+len(data))

	// 默认使用大端序
	binary.BigEndian.PutUint16(m, uint16(len(data)))

	copy(m[2:], data)

	// 发送消息
	n, err := conn.Write(m)
	log.Debug("n:%v, err:%v", n, err)

}
