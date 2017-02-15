package skeleton

import (
	"casino_common/common/log"
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/api"
	"sync"
	"casino_common/common/Error"
	"casino_common/common/consts"
	"casino_majiang/msg/funcsInit"
)

var ERR_SYS = Error.NewError(consts.ACK_RESULT_FAIL, "系统错误")
var ERR_REQ_REPETITION error = Error.NewError(consts.ACK_RESULT_FAIL, "重复请求")
var ERR_ENTER_DESK error = Error.NewError(consts.ACK_RESULT_FAIL, "进入房间失败")
var ERR_COIN_INSUFFICIENT error = Error.NewError(consts.ACK_RESULT_FAIL, "进入金币场失败，金币不足")

var ERR_LEAVE_RUNNING error = Error.NewError(consts.ACK_RESULT_FAIL, "现在不能离开")
var ERR_LEAVE_ERROR error = Error.NewError(consts.ACK_RESULT_FAIL, "出现错误，离开失败")

var ERR_READY_STATE = Error.NewError(consts.ACK_RESULT_FAIL, "准备失败，不在准备阶段")
var ERR_READY_REPETITION = Error.NewError(consts.ACK_RESULT_FAIL, "准备失败，不在准备阶段")
var ERR_READY_COIN_INSUFFICIENT = Error.NewError(consts.ACK_RESULT_FAIL, "准备失败，金币不足")
var ERR_READY_state = Error.NewError(consts.ACK_RESULT_FAIL, "准备失败，不在准备阶段")

//desk 的骨架,业务逻辑的方法 放置在这里
type SkeletonMJDesk struct {
	config   data.SkeletonMJConfig //这里不用使用指针，此配置创建之后不会再改变
	status   *data.MjDeskStatus    //桌子的所有状态都在这里
	HuParser api.HuPaerApi         //胡牌解析器
	*sync.Mutex
}

func NewSkeletonMJDesk(config data.SkeletonMJConfig) *SkeletonMJDesk {
	desk := &SkeletonMJDesk{
		config: config,
	}
	return desk
}

func (f *SkeletonMJDesk) EnterUser(userId uint32) error {
	log.Debug("玩家[%v]进入fdesk")
	return nil
}

//准备
func (d *SkeletonMJDesk) Ready(userId uint32) error {
	log.T("锁日志: %v ready(%v)的时候等待锁", d.DlogDes(), userId)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ready(%v)的时候释放锁", d.DlogDes(), userId)
	}()

	//判断desk状态
	if d.GetStatus().IsNotPreparing() {
		// 准备失败
		log.E("用户[%v]准备失败.desk[%v]不在准备的状态...", userId, d.GetMJConfig().DeskId)
		return ERR_READY_STATE
	}

	//找到需要准备的user
	user := d.GetUserByUserId(userId)
	if user == nil {
		log.E("用户[%v]在desk[%v]准备的时候失败,没有找到对应的玩家", userId, d.GetMJConfig().DeskId)
		return ERR_SYS
	}

	if user.GetStatus().IsReady() {
		log.E("玩家[%v]已经准备好了...请不要重新准备...", userId)
		return ERR_READY_REPETITION
	}

	//如果是金币场，需要判断玩家的金币是否足够
	//判断金币是否足够,准备的阶段不会扣除房费，房费是在开始的时候扣除

	user.Ready()

	//准备成功,发送准备成功的广播
	result := newProto.NewGame_AckReady()
	*result.Header.Code = consts.ACK_RESULT_SUCC
	*result.Header.Error = "准备成功"
	*result.UserId = userId
	log.T("广播user[%v]在desk[%v]准备成功的广播..string(%v)", userId, d.GetMJConfig().DeskId, result.String())
	d.BroadCastProto(result)

	return nil
}

//定缺
func (f *SkeletonMJDesk) DingQue(userId uint32, color int32) error {
	return nil
}
