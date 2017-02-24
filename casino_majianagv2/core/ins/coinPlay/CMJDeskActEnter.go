package coinPlay

import (
	"github.com/name5566/leaf/gate"
	"casino_common/proto/ddproto"
	"casino_majiang/msg/protogo"
	"casino_common/common/log"
	"casino_majiang/service/majiang"
	"casino_common/common/Error"
	"casino_majianagv2/core/ins/skeleton"
	"casino_common/common/consts"
)

//进入游戏
func (d *CMJDesk) EnterUser(userId uint32, a gate.Agent) error {
	log.T("锁日志: %v enterUser(%v)的时候等待锁", d.DlogDes(), userId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v enterUser(%v)的时候释放锁", d.DlogDes(), userId)
	}()

	rtype := d.TryReEnter(userId, a)
	if rtype {
		return nil
	}

	//先判断房间人数是否已经足够,这样就可以不处理下边的逻辑了
	if d.IsPlayerEnough() {
		return majiang.ERR_ENTER_DESK //进入房间失败
	}

	//3,加入一个新用户
	newUser := NewCMJUser(d, userId, a)
	if newUser == nil {
		return Error.NewError(consts.ACK_RESULT_ERROR, "创建user失败...")
	}

	snewUser := newUser.GetSkeletonUser().(*skeleton.SkeletonMJUser)
	if snewUser.Coin < d.CoinLimit {
		log.E("%v玩家进入房间失败,玩家[%v]的金币[%v]不足[%v]", d.DlogDes(), userId, snewUser.Coin, d.CoinLimit)
		return majiang.ERR_COIN_INSUFFICIENT
	}

	//扣除房费这里不用
	err := d.AddUserBean(newUser)
	if err != nil {
		log.E("用户[%v]加入房间[%v]失败,errMsg[%v]", userId, d.GetMJConfig().DeskId, err)
		return Error.NewFailError(Error.GetErrorMsg(err))
	} else {
		//加入房间成功,更新session  并且发送游戏数据
		snewUser.UpdateSession(int32(ddproto.COMMON_ENUM_GAMESTATUS_GAMING))
		d.SendGameInfo(userId, mjproto.RECONNECT_TYPE_NORMAL)
		//如果人数不够，这里会再初始化一个initEnterTimer,只有在金币场的时候会初始化
		return nil
	}
}
