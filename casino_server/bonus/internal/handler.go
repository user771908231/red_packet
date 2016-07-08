package internalimport (	"casino_server/msg/bbprotogo"	"reflect"	"github.com/name5566/leaf/gate"	"casino_server/common/log")//处理奖励相关的逻辑都写在这里//初始化func init(){	handler(&bbproto.LoginSignInBonus{}, handleLoginSignInBonus)	//17	登录签到奖励	handler(&bbproto.LoginTurntableBonus{}, handleLoginTurntableBonus)	//18	登录转盘奖励	handler(&bbproto.OlineBonus{}, handleOlineBonus)		//19	在线奖励	handler(&bbproto.TimingBonus{}, handleTimingBonus)		//20	定时奖励}func handler(m interface{}, h interface{}) {	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)}func handleLoginSignInBonus(args []interface{}){	log.T("进入了处理登录签到奖励的接口...handleLoginSignInBonus")	m := args[0].(*bbproto.Reg)	a := args[1].(gate.Agent)	log.T(m,a)}func handleLoginTurntableBonus(args []interface{}){	log.T("进入了处理登录转盘奖励的接口...handleLoginTurntableBonus")	m := args[0].(*bbproto.Reg)	a := args[1].(gate.Agent)	log.T(m,a)}func handleOlineBonus(args []interface{}){	log.T("进入了处理在线奖励的接口...handleOlineBonus")	m := args[0].(*bbproto.Reg)	a := args[1].(gate.Agent)	log.T(m,a)}func handleTimingBonus(args []interface{}){	log.T("进入了处理定时奖励的接口...handleTimingBonus")	m := args[0].(*bbproto.Reg)	a := args[1].(gate.Agent)	log.T(m,a)}