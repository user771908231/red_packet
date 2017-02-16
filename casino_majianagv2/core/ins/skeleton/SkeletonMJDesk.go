package skeleton

import (
	"casino_common/common/log"
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/api"
	"sync"
	"casino_common/common/Error"
	"casino_common/common/consts"
	"casino_majiang/msg/funcsInit"
	"github.com/name5566/leaf/timer"
	"casino_majiang/service/majiang"
	"casino_common/utils/rand"
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
	*sync.Mutex
	config        *data.SkeletonMJConfig //这里不用使用指针，此配置创建之后不会再改变
	status        *data.MjDeskStatus     //桌子的所有状态都在这里
	HuParser      api.HuPaerApi          //胡牌解析器
	OverTurnTimer *timer.Timer           //定时器
	CheckCase     *data.CheckCase        //麻将的判定器
	Users         []api.MjUser
}

func NewSkeletonMJDesk(config *data.SkeletonMJConfig) *SkeletonMJDesk {
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

//检测是否轮到当前玩家打牌...
func (d *SkeletonMJDesk) CheckActUser(userId uint32) bool {
	if d.GetMJConfig().ActUser == userId {
		return true //检测通过
	} else {
		//没有轮到当前玩家
		log.E("[%v]非法请求，没有轮到当前玩家打牌..应该是[%v]", d.DlogDes(), d.GetMJConfig().ActUser)
		return false
	}
}

//检测是否轮到操作
func (d *SkeletonMJDesk) CheckNotActUser(userId uint32) bool {
	return !d.CheckActUser(userId)
}

//定缺
func (f *SkeletonMJDesk) DingQue(userId uint32, color int32) error {
	return nil
}

func (d *SkeletonMJDesk) InitCheckCase(p *majiang.MJPai, outUser api.MjUser) error {
	return nil
}

//处理下一个checkCase
func (d *SkeletonMJDesk) DoCheckCase() error {
	return nil
}

//指针指向的玩家
func (d *SkeletonMJDesk) SetActiveUser(userId uint32) error {
	d.GetMJConfig().ActiveUser = userId
	return nil
}

func (d *SkeletonMJDesk) Time2Lottery() bool {
	//游戏中的玩家只剩下一个人，表示游戏结束...
	gamingCount := d.GetGamingCount() //正在游戏中的玩家数量

	log.T("判断是否胡牌...当前的gamingCount[%v],当前的PaiCursor[%v]", gamingCount, d.GetMJConfig().MJPaiCursor)
	//1,只剩下一个人的时候. 表示游戏结束
	if gamingCount == 1 {
		return true
	}

	log.T("HandPaiCanMo[%v]", d.HandPaiCanMo())
	//2，没有牌可以摸的时候，返回可以lottery了
	if !d.HandPaiCanMo() {
		return true
	}

	//如果是倒倒胡并且nextCheckCase为空
	if d.IsDaodaohu() && d.GetCheckCase().GetNextBean() == nil {
		for _, user := range d.Users {
			if user != nil && user.GetStatus().IsHu() {
				return true
			}
		}
	}
	return false
}

//发送game_opening 的协议
func (d *SkeletonMJDesk) SendNewGame_Opening() {
	log.T("发送游戏开始的协议..")
	log.T("当前桌子共[%v]把，现在是第[%v]把游戏开始", d.GetMJConfig().TotalPlayCount, d.GetMJConfig().CurrPlayCount)

	open := newProto.NewGame_Opening()
	*open.CurrPlayCount = d.GetMJConfig().CurrPlayCount
	*open.Dice1 = d.GetDice1() //骰子
	*open.Dice2 = d.GetDice2() //骰子
	d.BroadCastProto(open)
}

//通过庄来判断骰子的数目
func (d *SkeletonMJDesk) GetDice1() int32 {
	return rand.Rand(1, 7)
}

func (d *SkeletonMJDesk) GetDice2() int32 {
	return rand.Rand(1, 7)
}
