package console_command

import (
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"testing"
	"github.com/name5566/leaf/module"
	"github.com/name5566/leaf/chanrpc"
	"fmt"
)

type Module struct {
	*module.Skeleton
}

//初始化模块
func (m *Module) OnInit() {
	//骨架初始化
	skeleton := &module.Skeleton{
		GoLen:              10000,
		TimerDispatcherLen: 10000,
		ChanRPCServer:      chanrpc.NewServer(10000),
	}
	skeleton.Init()
	//配置初始化
	skeleton.RegisterCommand("echo", "echo test.", EchoCommand)

	m.Skeleton = skeleton
}

//模块销毁
func (m *Module) OnDestroy() {

}

//模块关闭
func (m *Module) Run(closeSig chan bool) {

}

//处理句柄
func EchoCommand(args []interface{}) interface{} {
	fmt.Printf("haha:%v", args)
	return fmt.Sprintf("haha:%v", args)
}

//测试：
func TestEchoCommanad(t *testing.T) {
	lconf.ConsolePort = 33888
	game_module := new(Module)
	leaf.Run(game_module)
}