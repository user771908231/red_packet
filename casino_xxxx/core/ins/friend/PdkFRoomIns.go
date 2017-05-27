package friend

import (
	"casino_paodekuai/core/ins/skeleton"
	"github.com/name5566/leaf/module"
	"casino_paodekuai/core/api"
	"casino_paodekuai/core/data"
	"casino_common/common/Error"
	"casino_common/common/consts"
	"casino_common/common/log"
)

type PdkFRoomIns struct {
	s *module.Skeleton        //eaf骨架
	*skeleton.PDKRoomSkeleton //room骨架
}

func NewPDKFRoom(s *module.Skeleton) *PdkFRoomIns {
	return &PdkFRoomIns{
		s:               s,
		PDKRoomSkeleton: skeleton.NewPDKRoomSkeleton(),
	}
}

var ERR_CREATE = Error.NewError(consts.ACK_RESULT_ERROR, "") //创建房间的错误

//创建房间
func (r *PdkFRoomIns) CreateDesk(cfg interface{}) (api.PDKDesk, error) {
	//处理创建房间的配置
	if cfg == nil {
		return nil, ERR_CREATE
	}
	deskCfg := cfg.(*data.PdkDeskCfg)

	//创建朋友桌的desk
	desk := NewPdkFDeskIns(deskCfg, r.s)
	if desk == nil {
		return nil, ERR_CREATE
	}

	//添加
	r.AddDesk(desk)

	log.T("创建房间的cfg:%v", deskCfg)
	return desk, nil
}

