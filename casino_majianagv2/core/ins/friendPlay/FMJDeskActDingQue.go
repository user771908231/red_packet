package friendPlay

import (
	"casino_common/common/log"
	"errors"
	"casino_majianagv2/core/majiangv2"
	"github.com/golang/protobuf/proto"
	"casino_majiang/msg/funcsInit"
	"casino_majiang/msg/protogo"
	"casino_majiang/service/majiang"
)

//个人开始定缺
func (d *FMJDesk) DingQue(userId uint32, color int32) error {
	log.T("锁日志: %v DingQue(%v,%v)的时候等待锁", d.DlogDes(), userId, color)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v DingQue(%v,%v)的时候释放锁", d.DlogDes(), userId, color, )
	}()

	//这里判断桌子的状态是否是在定缺的阶段
	if d.IsNotDingQue() {
		log.E("%v 当前的桌子不在定缺的状态", d.DlogDes())
		return errors.New("桌子当前不在定缺的状态，定缺失败")
	}
	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("定缺的时候，服务器出现错误，没有找到对应的user【%v】", userId)
		return errors.New("没有找到用户，定缺失败")
	}

	if user.GetStatus().DingQue {
		log.E("玩家[%v]重复定缺.", userId)
		return errors.New("用户已经定缺了，重复定缺....")
	}

	if color != majiangv2.T && color != majiangv2.S && color != majiangv2.W {
		log.E("玩家[%v]定缺的花色[%v]不合法.", userId, color)
		return errors.New("定缺失败.")
	}

	//设置定缺
	user.GetStatus().DingQue = true
	user.GetGameData().HandPai.QueFlower = proto.Int32(color)
	user.GetStatus().SetStatus(majiang.MJUSER_STATUS_DINGQUE) //设置目前的状态是已经定缺

	//回复定缺成功的消息
	ack := newProto.NewGame_DingQue()
	*ack.Header.UserId = userId
	*ack.Color = color
	*ack.UserId = userId
	log.T("广播：定缺成功的ack[%v]", ack)
	d.BroadCastProto(ack) //发送定缺成功的广播
	if d.AllDingQue() {
		//首先发送定缺结束的广播，然后发送庄家出牌的广播...
		log.T("%v 已经定缺完毕，现在看是发送定缺结束的协议...", d.DlogDes())
		d.GetDingQueEndInfo() //发送定缺完毕的广播
		d.BeginStart()        //游戏开始 庄家打牌

	}

	return nil
}

//定缺完毕
func (d *FMJDesk) GetDingQueEndInfo() *mjproto.Game_DingQueEnd {
	end := newProto.NewGame_DingQueEnd()
	for _, u := range d.GetUsers() {
		if u != nil && u.GetGameData().HandPai != nil {
			bean := newProto.NewDingQueEndBean()
			*bean.UserId = u.GetUserId()
			*bean.Flower = u.GetGameData().HandPai.GetQueFlower()
			end.Ques = append(end.Ques, bean)
		}
	}
	d.BroadCastProto(end)
	return end
}
