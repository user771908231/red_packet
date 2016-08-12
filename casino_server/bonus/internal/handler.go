package internal

import (
	"casino_server/msg/bbprotogo"
	"reflect"
	"github.com/name5566/leaf/gate"
	"casino_server/common/log"
	"casino_server/service/loginBonus"
)


//处理奖励相关的逻辑都写在这里


//初始化
func init(){
	handler(&bbproto.LoginSignInBonus{}, handleLoginSignInBonus)	//17	登录签到奖励
	handler(&bbproto.LoginTurntableBonus{}, handleLoginTurntableBonus)	//18	登录转盘奖励
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}



//处理登陆签到的奖励
func handleLoginSignInBonus(args []interface{}){
	log.T("进入了处理登录签到奖励的接口...handleLoginSignInBonus")
	m := args[0].(*bbproto.LoginSignInBonus)
	a := args[1].(gate.Agent)
	loginBonus.HandleLoginSignInBonus(m,a)
	log.T("",m,a)
}


//处理登录转盘的奖励

func handleLoginTurntableBonus(args []interface{}){
	log.T("进入了处理登录转盘奖励的接口...handleLoginTurntableBonus")
	m := args[0].(*bbproto.LoginTurntableBonus)
	a := args[1].(gate.Agent)
	loginBonus.HandleLoginTurntableBonus(m,a)
}

