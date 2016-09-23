package room

//游戏的主入口,启动两种不同的锦标赛都从这里进入

var DiamondGame *CSDiamondGame

func init() {
	DiamondGame = NewCsDiamondGame()
}

func StartGame() {
	DiamondGame.run(NewCsDiamondRoom())        //开始执行游戏
}