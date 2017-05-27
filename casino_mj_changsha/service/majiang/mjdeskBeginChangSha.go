package majiang

import (
	"casino_common/common/log"
	"github.com/golang/protobuf/proto"
	mjproto	"casino_mj_changsha/msg/protogo"
)

//开始游戏
/**
	开始游戏需要有几个步骤
	1，desk的状态是否正确，现在是否可以开始游戏

 */

func (d *MjDesk) beginChangSha() error {
	//因为begin是在准备之后才能准备，所以不用上锁，redy()上锁就可以了
	//1，检查是否可以开始游戏
	err := d.time2begin()
	if err != nil {
		log.T("无法开始游戏:err[%v]", err)
		return err
	}

	log.T("[%v]begin()...", d.DlogDes())

	//2，初始化桌子的状态
	d.beginInit()

	//3，根据playoptions 发cardsNum张牌
	err = d.initCards()
	if err != nil {
		log.E("初始化牌的时候出错err[%v]", err)
		return err
	}

	//判断起手胡牌
	err = d.checkChangShaQiShouHu()
	if err != nil {
		log.E("判断长沙牌起手胡的时候出错")
	}
	return nil
}

//检测长沙麻将的起手胡牌信息
func (d *MjDesk) checkChangShaQiShouHu() error {
	d.SetStatus(MJDESK_STATUS_QISHOUHU)
	//判断玩家是否可以起手胡牌，如果可以起手胡牌，那么发送起手胡牌的overTurn 给玩家
	//如果没有起手胡牌，那么直接开始begin
	err := d.initQishouhuCheckCase()
	if err != nil {
		log.E("%v 初始化起手胡的时候出错err%v", d.DlogDes(), err)
		return err
	}
	//如果有人可以起手胡牌，那么开始处理开始处理起手胡
	err = d.doQishouHuCheckCase()
	if err != nil {
		log.E("%v 开始处理起手胡牌的时候出错..")
		return err
	}
	return nil
}

//初始化起手胡的信息
func (d *MjDesk) initQishouhuCheckCase() error {
	d.CheckCase = NewCheckCase()
	d.CheckCase.DianPaoCount = proto.Int32(0)
	for _, u := range d.GetUsers() {
		if u != nil {
			hu, _, _, _, _, _ := d.HuParser.GetCanHu(u.GameData.GetHandPai(), nil, false, mjproto.HuType_H_changsha_qishouhu, d.IsBanker(u))
			if hu {
				checkBean := &CheckBean{}
				checkBean.CanHu = proto.Bool(true)
				checkBean.UserId = proto.Uint32(u.GetUserId())
				checkBean.CheckStatus = proto.Int32(CHECK_CASE_BEAN_STATUS_CHECKING)
				d.CheckCase.CheckB = append(d.CheckCase.CheckB, checkBean)
			}
		}
	}

	//判断是否有需要判断的checkBean
	if d.CheckCase.GetCheckB() == nil || len(d.CheckCase.GetCheckB()) == 0 {
		d.CheckCase = nil
	}
	return nil
}

//处理玩家的起手胡overTurn
func (d *MjDesk) doQishouHuCheckCase() error {
	caseBean := d.CheckCase.GetNextBean()
	d.UpdateUserStatus(MJUSER_STATUS_GAMING) //设置为游戏中的状态
	//检测参数
	if caseBean == nil {
		d.CheckCase = nil //设置CheckCase为nil
		d.BeginStart()    //没有起手胡，开始游戏
		return nil
	} else {
		log.T("继续处理CheckCase,开处理下一个checkBean...")
		//找到需要判断bean之后，发送给判断人	//send overTurn
		overTurn := &mjproto.Game_ChangshQiShouHuOverTurn{
			Header: &mjproto.ProtoHeader{
				UserId: caseBean.UserId,
			},
		}
		log.T("开始发送overTurn[%v]", overTurn)
		d.GetUserByUserId(caseBean.GetUserId()).SendOverTurn(overTurn)        //DoCheckCase
		d.SetActUserAndType(caseBean.GetUserId(), MJDESK_ACT_TYPE_WAIT_CHECK) // 起手胡牌 DoCheckCase 设置当前活动的玩家
		return nil
	}
}
