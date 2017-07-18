package protobuf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/log"
	clog "casino_common/common/log"
	"math"
	"reflect"
	"time"
)

// -------------------------
// | id | protobuf message |
// -------------------------
type Processor struct {
	littleEndian bool
	msgInfo      []*MsgInfo
	msgID        map[reflect.Type]uint16
}

type MsgInfo struct {
	msgType    reflect.Type
	msgRouter  *chanrpc.Server
	msgHandler MsgHandler
}

type MsgHandler func([]interface{})

func NewProcessor() *Processor {
	p := new(Processor)
	p.littleEndian = false
	p.msgID = make(map[reflect.Type]uint16)
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) Register(msg proto.Message) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("protobuf message pointer required")
	}
	if _, ok := p.msgID[msgType]; ok {
		log.Fatal("message %s is already registered", msgType)
	}
	if len(p.msgInfo) >= math.MaxUint16 {
		log.Fatal("too many protobuf messages (max = %v)", math.MaxUint16)
	}

	i := new(MsgInfo)
	i.msgType = msgType
	p.msgInfo = append(p.msgInfo, i)
	p.msgID[msgType] = uint16(len(p.msgInfo) - 1)
	//log.Debug("注册: p.msgID[msgType]    %v ",p.msgID[msgType])
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRouter(msg proto.Message, msgRouter *chanrpc.Server) {
	msgType := reflect.TypeOf(msg)
	id, ok := p.msgID[msgType]
	if !ok {
		log.Fatal("message %s not registered", msgType)
	}

	p.msgInfo[id].msgRouter = msgRouter
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetHandler(msg proto.Message, msgHandler MsgHandler) {
	msgType := reflect.TypeOf(msg)
	id, ok := p.msgID[msgType]
	if !ok {
		log.Fatal("message %s not registered", msgType)
	}

	p.msgInfo[id].msgHandler = msgHandler
}

// goroutine safe
func (p *Processor) Route(msg interface{}, userData interface{}) error {
	msgType := reflect.TypeOf(msg)
	id, ok := p.msgID[msgType]
	if !ok {
		return fmt.Errorf("message %s not registered", msgType)
	}

	i := p.msgInfo[id]

	if i.msgHandler != nil {
		//log.Debug("Processor.Route(): 执行msgHandler()：%v", msg)  //TODO: 临时调试log
		time_start := time.Now()
		i.msgHandler([]interface{}{msg, userData})
		time_spend := time.Now().Sub(time_start).Seconds() * 1e3
		log.Debug("Processor.Route(): 执行msgHandler()：%T spend:[%.2f ms]", msg, time_spend)  //TODO: 临时调试log
		if time_spend >= 100 {
			clog.E("Processor.Route(): 执行msgHandler()：%T spend:[%.2f ms]", msg, time_spend)
		}
	}

	if i.msgRouter != nil {
		//log.Debug("Processor.Route(): 执行msgRouter.Go()：%v", msg)  //TODO: 临时调试log
		i.msgRouter.Go(msgType, msg, userData)
	}
	return nil
}

// goroutine safe
func (p *Processor) Unmarshal(data []byte) (interface{}, error) {

	if len(data) < 2 {
		return nil, errors.New("protobuf data too short")
	}

	// id
	var id uint16
	//log.Debug("p.littleEndian: %v",p.littleEndian)
	if p.littleEndian {
		id = binary.LittleEndian.Uint16(data)
	} else {
		id = binary.BigEndian.Uint16(data)
	}
	//log.Debug("protobuf 格式的id: %v",id)
	// msg
	if id >= uint16(len(p.msgInfo)) {
		return nil, fmt.Errorf("message id %v not registered", id)
	}
	msg := reflect.New(p.msgInfo[id].msgType.Elem()).Interface()

	//log.Debug("protobuf 的内容 data[2:]: %v",data[2:])

	return msg, proto.UnmarshalMerge(data[2:], msg.(proto.Message))
}

// goroutine safe
func (p *Processor) Marshal(msg interface{}) ([][]byte, error) {
	msgType := reflect.TypeOf(msg)

	// id
	_id, ok := p.msgID[msgType]
	if !ok {
		err := fmt.Errorf("message %s not registered", msgType)
		return nil, err
	}

	id := make([]byte, 2)
	if p.littleEndian {
		binary.LittleEndian.PutUint16(id, _id)
	} else {
		binary.BigEndian.PutUint16(id, _id)
	}

	// data
	data, err := proto.Marshal(msg.(proto.Message))
	return [][]byte{id, data}, err
}

// goroutine safe
func (p *Processor) Range(f func(id uint16, t reflect.Type)) {
	for id, i := range p.msgInfo {
		f(uint16(id), i.msgType)
	}
}
