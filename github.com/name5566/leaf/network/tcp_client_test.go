package network

//
//import (
//	"testing"
//	"leaf-note/leaf/network"
//	"time"
//	"math"
//)
//
//func TestTCPClient(t *testing.T) {
//
//	var addr string = "192.168.199.120:2801"
//	var PendingWriteNum int = 1
//	client := new(network.TCPClient)
//	client.Addr = addr
//	client.ConnNum = 1
//	client.ConnectInterval = 3 * time.Second
//	client.PendingWriteNum = PendingWriteNum
//	client.LenMsgLen = 4
//	client.MaxMsgLen = math.MaxUint32
//	client.NewAgent = newAgent
//	client.Start()
//}
//
//type MAgent struct {
//	conn *network.TCPConn
//}
//
//func newAgent(conn *network.TCPConn) network.Agent {
//	a := new(MAgent)
//	a.conn = conn
//	return a
//}
//
//func (a *MAgent) Run() {}
//
//func (a *MAgent) OnClose() {}
