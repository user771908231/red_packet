package friendPlay

import (
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
	"casino_majiang/service/majiang"
	"casino_common/common/log"
)

//离开房间
func (d *FMJDesk) Leave(userId uint32) error {

	//1
	if d.GetStatus().IsGaming() {
		return majiang.ERR_LEAVE_RUNNING //离开房间失败
	}

	//2,准备阶段的时候可以离开
	if d.GetStatus().IsNotPreparing() {
		return majiang.ERR_LEAVE_RUNNING //只有在准备的时候才可以离开
	}

	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("玩家[%v]离开的时候失败，没有找到对应的user", user.GetUserId())
		return majiang.ERR_LEAVE_ERROR
	}

	err := d.RmUser(user)
	if err != nil {
		return err
	}

	//发送离开的广播
	ack := new(ddproto.CommonAckLeaveDesk)
	ack.UserId = proto.Uint32(user.GetUserId())
	ack.IsExchange = proto.Bool(false)
	d.BroadCastProto(ack)

	if d.canDissolve() {
		dissolveErr := d.Room.DissolveDesk(d, false)
		if dissolveErr != nil {
			log.E("解散房间的时候失败%v", dissolveErr)
			return nil
		}

	}
	return nil
}

// 是否满足解散的条件
func (d *FMJDesk) canDissolve() bool {
	//1,申请解散的人是否达标
	userCount := d.GetUserCount()
	log.T("%v 目前总人数:%v", d.DlogDes(), userCount)

	//2,如果user为0
	if userCount == 0 {
		return true
	}

	//3,检测申请解散的
	applyCount := int32(len(d.AllApplyDissolveUser()))
	log.T("%v 目前同意解散的人数:%v", d.DlogDes(), applyCount)
	if applyCount == userCount {
		return true
	}
	return false
}

func (d *FMJDesk) AllApplyDissolveUser() []uint32 {
	ret := make([]uint32, 0)
	for _, u := range d.GetUsers() {
		if u != nil && u.GetStatus().GetApplyDissolve() == majiang.MJUER_APPLYDISSOLVE_S_AGREE {
			log.T("%v检测用户同意解散了", d.DlogDes(), u.GetUserId())
			ret = append(ret, u.GetUserId())
		}
	}
	return ret
}
