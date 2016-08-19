package room

import (
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"time"
	"casino_server/utils/db"
	"casino_server/conf/casinoConf"
	"casino_server/mode"
	"gopkg.in/mgo.v2/bson"
	"errors"
)

var ChampionshipRoom CSThGameRoom 	//锦标赛的房间


func init(){
	ChampionshipRoom.OnInit()	//初始化,开始运行
	ChampionshipRoom.Run()
}
//锦标赛
type CSThGameRoom struct {
	ThGameRoom
	//锦标赛房间的专有属性
	matchId	int32		//比赛内容
	beginTime	time.Time	//游戏开始的时间
	endTime 	time.Time	//游戏结束的时间
}

//只有开始之后才能进入游戏房间
func (r *CSThGameRoom) IsBegin() bool{
	nowTime := time.Now()
	if nowTime.Before(r.endTime) {
		return true
	}else{
		return false
	}
}


//run游戏房间
func (r *CSThGameRoom) Run() {
	log.T("锦标赛游戏开始...")

	//设置room属性
	r.beginTime = time.Now()
	r.endTime   = r.beginTime.Add(time.Second*60*20)		//一局游戏的时间是20分钟
	r.matchId   = db.GetNextSeq(casinoConf.DBT_T_CS_TH_RECORD)	//生成游戏的matchId

	//保存游戏数据
	saveData := &mode.T_cs_th_record{}
	saveData.Mid = bson.NewObjectId()
	saveData.Id = r.matchId
	saveData.BeginTime = r.beginTime
	saveData.EndTime = r.endTime
	db.SaveMgoData(casinoConf.DBT_T_CS_TH_RECORD,saveData)


	//这里定义一个计时器,每十秒钟检测一次游戏
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for timeNow := range ticker.C {
			log.T("开始time[%v]检测锦标赛matchId[%v]有没有结束...",timeNow,r.matchId)
			if r.checkEnd() {
				//重新开始
				go r.Run()
				break
			}
		}
	}()
}

//检测结束
func (r *CSThGameRoom) checkEnd() bool{
	//如果时间已经过了,并且所有桌子的状态都是已经停止游戏,那么表示这一局结束
	if r.endTime.Before(time.Now()) && r.allStop() {
		//结算本局
		return true
	}else{
		return  false
	}

}


//判断是否所有的desk停止游戏
func (r *CSThGameRoom) allStop() bool{
	result := true
	for i := 0; i < len(r.ThDeskBuf); i++ {
		desk := r.ThDeskBuf[i]
		if  desk != nil && desk.Status != TH_DESK_STATUS_STOP{
			result = false
			break
		}
	}
	return result

}

func (r *CSThGameRoom) End(){
	log.T("锦标赛游戏结束")
}

//游戏大厅增加一个玩家
func (r *CSThGameRoom) AddUser(userId uint32, roomCoin int64, a gate.Agent) (*ThDesk, error) {
	r.Lock()
	defer r.Unlock()
	log.T("userid【%v】进入德州扑克的房间", userId)

	//这里需要判断锦标赛是否可以开始游戏
	if !r.IsBegin() {
		log.T("用户[%v]进入锦标赛的房间失败,因为游戏还没有开始",userId)
		return nil,errors.New("游戏还没有开始")
	}

	var mydesk *ThDesk = nil                //为用户找到的desk
	//1,判断用户是否已经在房间里了,如果是在房间里,那么替换现有的agent,
	mydesk = r.IsRepeatIntoRoom(userId, a)
	if mydesk != nil {
		return mydesk, nil
	}

	//2,查询哪个德州的房间缺人:循环每个德州的房间,然后查询哪个房间缺人
	for deskIndex := 0; deskIndex < len(r.ThDeskBuf); deskIndex++ {
		tempDesk := r.ThDeskBuf[deskIndex]
		if tempDesk == nil {
			log.E("找到房间为nil,出错")
			break
		}
		if tempDesk.UserCount < r.ThRoomSeatMax {
			mydesk = tempDesk        //通过roomId找到德州的room
			break;
		}
	}

	//如果没有可以使用的桌子,那么重新创建一个,并且放进游戏大厅
	if mydesk == nil {
		log.T("没有多余的desk可以用,重新创建一个desk")
		mydesk = NewThDesk()
		mydesk.MatchId = r.matchId
		r.AddThDesk(mydesk)
	}

	//3,进入房间,竞标赛进入房间的时候,默认就是准备的状态
	err := mydesk.AddThUser(userId, roomCoin, TH_USER_STATUS_READY, a)
	if err != nil {
		log.E("用户上德州扑克的桌子 失败...")
		return nil, err
	}

	mydesk.LogString()        //答应当前房间的信息
	return mydesk, nil
}


//是否可以进行下把游戏
func (r *CSThGameRoom) CanNextDeskRun() bool{
	nowTime := time.Now()
	if r.endTime.Before(nowTime) {
		//如果当前时间已经在结束时间之后,那么本局游戏结束
		return false
	}
	return true
}

