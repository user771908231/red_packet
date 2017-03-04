package coinPlay

import (
	"casino_majiang/service/majiang"
	"casino_majianagv2/core/majiangv2"
	"casino_common/common/Error"
	"casino_common/common/log"
	"casino_common/common/consts"
	"casino_majiang/msg/funcsInit"
	"errors"
	"github.com/name5566/leaf/util"
	"casino_majiang/msg/protogo"
	"github.com/golang/protobuf/proto"
)

var ERR_OUTPAI error = Error.NewError(consts.ACK_RESULT_ERROR, "打牌失败")

//打牌
func (d *CMJDesk) ActOut(userId uint32, paiKey int32, auto bool) error {
	defer Error.ErrorRecovery("actOut")
	log.T("锁日志: %v ActOut(%v,%v)的时候等待锁", d.DlogDes(), userId, paiKey)
	d.Lock()
	defer func() {
		d.Unlock()
		log.T("锁日志: %v ActOut(%v,%v)的时候释放锁", d.DlogDes(), userId, paiKey, )
	}()

	outUser := d.GetUserByUserId(userId)
	if outUser == nil {
		log.E("[%v]打牌失败，没有找到玩家[%v]", d.DlogDes(), userId)
		return ERR_OUTPAI
	}

	//判断是否轮到当前玩家打牌了...
	if d.CheckNotActUser(userId) { //打牌
		log.E("%v没有轮到当前玩家%v打牌", d.DlogDes(), userId)
		return ERR_OUTPAI
	}

	//判断是不是在游戏中的状态
	if d.GetStatus().IsNotGaming() {
		log.E("%v玩家%v打牌失败，desk不在游戏状态[%v]", d.DlogDes(), userId, d.GetStatus().IsNotGaming())
		return ERR_OUTPAI
	}

	//停止定时器
	if d.OverTurnTimer != nil {
		d.OverTurnTimer.Stop()
	}

	//得到参数
	outPai := majiangv2.InitMjPaiByIndex(int(paiKey))

	/**
		1,如果是碰牌打牌的时候,inpai为nil，不需要增加
		2,如果是摸牌打牌（杠之后也是摸牌，需要增加in牌...）
	 */

	if outUser.GetGameData().HandPai.InPai != nil {
		outUser.GetGameData().HandPai.AddPai(outUser.GetGameData().HandPai.InPai) //把inpai放置到手牌上
	}
	log.T("玩家打牌之前的手牌:%v", majiangv2.ServerPais2string(outUser.GetGameData().HandPai.Pais))
	errDapai := outUser.GetGameData().HandPai.DelHandlPai(outPai.GetIndex()) //删除要打出去的牌
	if errDapai != nil {
		log.E("[%v]打牌的时候出现错误，没有找到要到的牌,id[%v]", d.DlogDes(), paiKey)
		return ERR_OUTPAI
	}
	log.T("玩家打牌之后的手牌:%v", majiangv2.ServerPais2string(outUser.GetGameData().HandPai.Pais))

	outUser.GetGameData().HandPai.OutPais = append(outUser.GetGameData().HandPai.OutPais, outPai) //自己桌子前面打出的牌，如果其他人碰杠胡了之后，需要把牌删除掉...
	outUser.GetGameData().HandPai.InPai = nil                                                     //打牌之后需要把自己的  inpai给移除掉...
	outUser.GetGameData().DelGuoHuInfo()                                                          //删除过胡的信息
	//打牌之后的逻辑,初始化判定事件
	err := d.InitCheckCase(outPai, outUser) //打牌之后
	if err != nil {
		//表示无人需要，直接给用户返回无人需要
		//给下一个人摸排，并且移动指针
		log.E("%v服务器错误，初始化判定牌的时候出错err[%v]", d.DlogDes(), err)
		return ERR_OUTPAI
	}

	log.T("[%v]玩家[%v]打牌之后InitCheckCase之后的checkCase[%v]", d.DlogDes(), userId, d.CheckCase)
	//回复消息,打牌之后，广播打牌的信息...s
	outUser.GetGameData().DelPreMoGangInfo() //清楚摸牌前的杠牌info
	result := newProto.NewGame_AckSendOutCard()
	*result.UserId = userId
	result.Card = outPai.GetCardInfo()
	d.BroadCastProto(result)

	log.T("[%v]用户[%v]已经打牌结束，开始处理下一个checkCase", d.DlogDes(), userId)
	d.DoCheckCase() //打牌之后，别人判定牌
	return nil
}

func (d *CMJDesk) DoCheckCase() error {
	//检测参数
	if d.CheckCase.GetNextBean() == nil {
		log.T("[%v]已经没有需要处理的CheckCase,下一个玩家摸牌...", d.DlogDes())
		//直接跳转到下一个操作的玩家...,这里表示判断已经玩了...
		d.CheckCase = nil
		//在这之前需要保证 activeUser 是正确的...
		d.SendMopaiOverTurnChengDu()
		return nil
	} else {
		log.T("继续处理CheckCase,开处理下一个checkBean...")
		//1,找到胡牌的人来进行处理
		caseBean := d.CheckCase.GetNextBean()
		//找到需要判断bean之后，发送给判断人	//send overTurn
		overTurn := d.GetOverTurnByCaseBean(d.CheckCase.CheckMJPai, caseBean, majiang.OVER_TURN_ACTTYPE_OTHER) //别人打牌，判断是否可以碰杠胡

		///发送overTurn 的信息
		log.T("%v 开始发送overTurn[%v]", d.DlogDes(), overTurn)
		d.GetUserByUserId(caseBean.GetUserId()).SendOverTurn(overTurn)                //DoCheckCase
		d.SetActUserAndType(caseBean.GetUserId(), majiang.MJDESK_ACT_TYPE_WAIT_CHECK) //DoCheckCase 设置当前活动的玩家
		return nil
	}
}

//发送摸牌的广播
//指定一个摸牌，如果没有指定，则系统通过游标来判断
func (d *CMJDesk) SendMopaiOverTurnChengDu() error {
	//首先判断是否可以lottery(),如果可以那么直接开奖
	if d.Time2Lottery() {
		d.LotteryChengDu() //摸牌的时候判断可以lottery了
		return nil
	}

	//开始摸牌的逻辑
	user := d.GetNextMoPaiUser()
	if user == nil {
		log.E("服务器出现错误..没有找到下一个摸牌的玩家...")
		return errors.New("没有找到下一家")
	}

	//转换user

	d.SetActiveUser(user.GetUserId())                                    //用户摸牌之后，设置前端指针指向的玩家
	d.SetActUserAndType(user.GetUserId(), majiang.MJDESK_ACT_TYPE_MOPAI) //长度麻将 用户摸牌之后，设置当前活动的玩家
	//发送摸牌的OverTrun
	user.GetGameData().GetHandPai().InPai = d.GetNextPai()
	overTrun := d.GetMoPaiOverTurn(user, false) //普通摸牌，用户摸牌的时候,发送一个用户摸牌的overturn
	user.SendOverTurn(overTrun)                 //玩家摸排之后发送overturn

	//给其他人广播协议
	overTurn2 := &mjproto.Game_OverTurn{}
	util.DeepCopy(overTurn2, overTrun)
	overTurn2.CanHu = proto.Bool(false)
	overTurn2.CanGang = proto.Bool(false)
	overTurn2.ActCard = majiang.NewBackPai()
	d.BroadCastProtoExclusive(overTrun, user.GetUserId())
	return nil
}
