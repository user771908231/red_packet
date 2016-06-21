package fruitService

import (
	"casino_server/msg/bbproto"
	"github.com/name5566/leaf/gate"
	"casino_server/conf/intCons"
)


/**
	得到一次的结果
	水果机器的结果有可能有很多种,这里需要什么策略来返回结果?
 */


func HandlerShuiguoji(m *bbproto.Shuiguoji,a gate.Agent) (*bbproto.Shuiguoji,error){
	result := &bbproto.Shuiguoji{

	}

	//检测参数并且根据押注的内容选择处理方式
	if m == nil {
		return result,nil
	}

	vtype := m.GetType()
	switch vtype {
	case intCons.SGJV.TypeBet:
		BetResult(m.GetProtoHeader().GetUserId())
	case intCons.SGJV.TypeHilomp:
		HilompResult(m.GetProtoHeader().GetUserId())
	}


	return result,nil

}

func BetResult(id uint32) error{
	return nil
}

/**
	比大小的结果
 */
func HilompResult(id uint32) error{
	return nil

}
