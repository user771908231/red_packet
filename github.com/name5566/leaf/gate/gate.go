package gate

import (
	clog "casino_common/common/log"
	"casino_common/utils/security"
	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/network"
	"reflect"
	"strings"
	"time"
	"casino_common/common/Error"
	"fmt"
)

func init() {
	agentList = map[*agent]bool{}
}

type Gate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       network.Processor
	AgentChanRPC    *chanrpc.Server

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool
}

//连接上次接收消息的时间
var agentList map[*agent]bool

func (gate *Gate) Run(closeSig chan bool) {
	var wsServer *network.WSServer
	if gate.WSAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = gate.WSAddr
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.PendingWriteNum = gate.PendingWriteNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			a.lastReceiveTime = time.Now()
			agentList[a] = true
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go("NewAgent", a)
			}
			return a
		}
	}

	var tcpServer *network.TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go("NewAgent", a)
			}
			return a
		}
	}

	if wsServer != nil {
		wsServer.Start()
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	//自动销毁超时的agent
	go func() {
		for {
			//每次循环延时60秒
			<-time.After(60 * time.Second)

			time_now := time.Now()
			log.Debug("start clean timeout agent, curr agentlist [len: %d]", len(agentList))
			for a, _ := range agentList {
				if a == nil {
					continue
				}
				timeCost := time_now.Sub(a.lastReceiveTime)
				if timeCost > time.Second * 60 {
					log.Debug("agent:%p close [timeout: %.2f s]", a, timeCost.Seconds())
					//超时则关闭链接，并从列表中删除
					a.Close()
					a.Destroy()
					delete(agentList, a)
				}
			}
			log.Debug("end clean timeout agent, curr agentlist [len: %d]", len(agentList))
		}
	}()
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (gate *Gate) OnDestroy() {}

type agent struct {
	conn     network.Conn
	gate     *Gate
	userData interface{}
	lastReceiveTime time.Time
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()

		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		//更新上次收到数据的时间
		a.lastReceiveTime = time.Now()

		if a.gate.Processor != nil {
			//增加一层校验md5的方法
			data2, checkErr := security.CheckTcpData(data)
			if checkErr != nil {
				log.Debug("data check md5 fail: %v", checkErr)
				break
			}

			//data2 := data
			msg, err := a.gate.Processor.Unmarshal(data2)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}

			//打印接收到的信息
			typeString := reflect.TypeOf(msg).String()
			log.Debug("agent[%p]解析出来的数据type[%v],m[%v]", a, reflect.TypeOf(msg).String(), msg)
			if !strings.Contains(typeString, "Heartbeat") {
				clog.T("a[%p]解析出来的数据type[%v],m[%v]", a, reflect.TypeOf(msg).String(), msg)
			}

			err = func(err error) error {
				defer Error.ErrorRecovery(fmt.Sprintf("a.gate.Processor.Route(%v, %p)", typeString, a))
				return a.gate.Processor.Route(msg, a)
			}(err)
			if err != nil {
				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (a *agent) OnClose() {
	if a.gate.AgentChanRPC != nil {
		err := a.gate.AgentChanRPC.Open(0).Call0("CloseAgent", a)
		if err != nil {
			log.Error("chanrpc error: %v", err)
		}
	}
}

func (a *agent) WriteMsg(msg interface{}) {
	typeString := reflect.TypeOf(msg).String()
	if a.gate.Processor != nil {
		data, err := a.gate.Processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		time_start := time.Now()
		a.conn.WriteMsg(data...)
		time_end := time.Now()
		time_sub := time_end.Sub(time_start)
		log.Debug("agent[%p]发送的信息 type[%v],id[%v] len[%v] spend[%.2f ms],\t\t content[%v]", a, typeString, data[0], len(data[1]), time_sub.Seconds()*1e3, msg)
	}
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}

func (a *agent) RemoteAddr() interface{} {
	return a.conn.RemoteAddr()
}
