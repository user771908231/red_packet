package majiang

import (
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
	"casino_common/common/log"
	"casino_common/common/Error"
)

// 是否满足解散的条件
func (d *MjDesk) canDissolve() bool {
	//return false
	//朋友桌的条件
	if d.IsFriend() {
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

	}

	return false
}

func (d *MjDesk) AllApplyDissolveUser() []uint32 {
	ret := make([]uint32, 0)
	for _, u := range d.GetUsers() {
		if u != nil && u.GetApplyDissolve() == MJUER_APPLYDISSOLVE_S_AGREE {
			log.T("%v检测用户同意解散了", d.DlogDes(), u.GetUserId())
			ret = append(ret, u.GetUserId())
		}
	}
	return ret
}

func (d *MjDesk) getOnlineCount() int32 {
	ret := int32(0)
	for _, u := range d.GetUsers() {
		if u != nil && !u.GetIsBreak() && !u.GetIsLeave() {
			ret++
		}
	}
	return ret
}

func (d *MjDesk) AllUnApplyDissolveUser() []uint32 {
	ret := make([]uint32, 0)
	for _, u := range d.GetUsers() {
		if u != nil && u.GetApplyDissolve() == MJUER_APPLYDISSOLVE_S_REFUSE {
			ret = append(ret, u.GetUserId())
		}
	}
	return ret
}

//申请解散房间
func (d *MjDesk) ApplyDissolve(userId uint32) error {
	log.T("锁日志: %v ApplyDissolve(%v)的时候等待锁", d.DlogDes(), userId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ApplyDissolve(%v)的时候释放锁", d.DlogDes(), userId)
	}()

	d.ApplyDis = proto.Bool(true)     //有人申请解散房间:设置为等待解散的状态
	user := d.GetUserByUserId(userId) //得到申请人

	user.applayDissolve(MJUER_APPLYDISSOLVE_S_AGREE)
	//3,如果可以解散 直接解散房间
	if d.canDissolve() {
		//达到解散的条件
		err := MjroomManagerIns.GetFMJRoom().DissolveDesk(d, true) //申请解散房间，开始解散房间
		if err != nil {
			log.E("解散房间的时候发生错误:%v", err)
		}
	} else {
		//4,如果房间还没有解散，那么现在开始做定时处理
		for _, u := range d.GetUsers() {
			if u != nil {
				userId := u.GetUserId()
				//短线离线的人 默认同意.其他的人增加定时器
				if u.GetIsBreak() || u.GetIsLeave() {
					//u.applayDissolve(MJUER_APPLYDISSOLVE_S_AGREE) //设置申请解散状态
					go func() {
						defer Error.ErrorRecovery("离线的人自动统一")
						d.ApplyDissolveBack(userId, true)
					}()
				} else if u.GetUserId() != user.GetUserId() && u.GetApplyDissolve() == MJUER_APPLYDISSOLVE_S_DEFAULT {
					log.T("开始给user %v 设置同意解散房间倒计时", userId)
					u.dissolveTimer = d.AfterFunc(APPLYDISSOLVE_DURATION, func() {
						defer Error.ErrorRecovery("倒计时解散房间")
						//超时的时候，需要默认设置为同意
						log.T("玩家[%v]同意或者拒绝解散房间的倒计时超时...现在开始进行超时处理.", userId)
						d.ApplyDissolveBack(userId, true)
					})
				}
			}
		}
	}

	//2,申请解散 回复
	ack := new(ddproto.CommonBcApplyDissolve)
	ack.UserId = proto.Uint32(userId)
	d.BroadCastProto(ack)

	return nil
}

//申请解散房间
func (d *MjDesk) ApplyDissolveBack(userId uint32, agree bool) error {
	log.T("锁日志: %v ApplyDissolveBack(%v,%v)的时候等待锁", d.DlogDes(), userId, agree)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ApplyDissolveBack(%v,%v)的时候释放锁", d.DlogDes(), userId, agree)
	}()

	log.T("ApplyDissolveBack(%v,%v),,d.GetApplyDis():%v", userId, agree, d.GetApplyDis())

	if !d.GetApplyDis() {
		//不在 判断拒绝或者同意的状态
		//return ERR_REQ_REPETITION
		log.W("玩家重复请求解散房间")
		return nil
	}

	//如果有人拒绝则表示申请解散失败...
	if !agree {
		log.T("%v 由于玩家[%v]agree:%v 所以解散失败..", d.DlogDes(), userId, agree)
		//解散失败，需要把其他人申请的状态设置为 没有统一
		for _, u := range d.GetUsers() {
			if u != nil {
				u.applayDissolve(MJUER_APPLYDISSOLVE_S_DEFAULT) //有人拒绝以后设置为默认状态
				if u.dissolveTimer != nil {
					u.dissolveTimer.Stop() //停止定时器
				}
			}
		}
		d.ApplyDis = proto.Bool(false) //有人不同意解散房间：设置不在申请的阶段...
		ack := new(ddproto.CommonAckApplyDissolveBack)
		ack.UserId = proto.Uint32(userId)
		d.BroadCastProto(ack)
		return nil
	}

	//得到用户，并判断是否是重复请求，如果不是重复请求就设置状态，并返回ack
	user := d.GetUserByUserId(userId)
	if user.GetApplyDissolve() != MJUER_APPLYDISSOLVE_S_DEFAULT {
		return ERR_REQ_REPETITION
	}
	//设置等待状态
	user.applayDissolve(MJUER_APPLYDISSOLVE_S_AGREE) //设置申请解散状态
	if user.dissolveTimer != nil {
		user.dissolveTimer.Stop() //停止定时器
	}

	//回复解散的信息
	ack := new(ddproto.CommonAckApplyDissolveBack) //回复ack
	ack.UserId = proto.Uint32(userId)
	ack.Agree = proto.Bool(agree)
	d.BroadCastProto(ack)

	//如果可以解散 直接解散房间
	if d.canDissolve() {
		d.ApplyDis = proto.Bool(false)                             //可以解散房间，设置为不在申请的状态
		err := MjroomManagerIns.GetFMJRoom().DissolveDesk(d, true) //统一解散房间 开始解散房间
		if err != nil {
			log.E("解散房间的时候发生错误:%v", err)
		}
		return err
	}

	return nil
}
