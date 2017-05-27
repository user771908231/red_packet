package friend

import (
	"github.com/name5566/leaf/module"
	"casino_paodekuai/core/ins/skeleton"
	"casino_paodekuai/core/data"
	"casino_common/proto/ddproto"
	"src/github.com/golang/protobuf/proto"
	"casino_common/common/consts"
)

//跑得快朋友桌
type PdkFDeskIns struct {
	s   *module.Skeleton
	*skeleton.PDKDeskSkeleton
	cfg *data.PdkDeskCfg
}

func NewPdkFDeskIns(cfg *data.PdkDeskCfg, s *module.Skeleton) *PdkFDeskIns {
	return &PdkFDeskIns{
		s:               s,
		PDKDeskSkeleton: skeleton.NewPDKDeskSkeleton(),
	}
}

//todo 进入房间
func (d *PdkFDeskIns) EnterUser(userId uint32) error {

	ack := &ddproto.PdkAckEnterDesk{
		Header:&ddproto.ProtoHeader{
			UserId:proto.Uint32(userId),
			Code:proto.Int32(consts.ACK_RESULT_SUCC),
		},
	}
	d.BC(ack)
	return nil
}

//出牌的user和牌型
func (d *PdkFDeskIns) ActOut(userId uint32, p interface{}) error {
	//1,p 要打的牌
	//2，判断p 是不是比一家的牌大
	//3, user.outCard()玩家 把手里面的牌打出去
	//4, 通知其他玩家打的啥牌


	return nil
}
