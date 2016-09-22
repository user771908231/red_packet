package room

//游戏的主入口,启动两种不同的锦标赛都从这里进入
func StartGame() {
	NewCsDiamondGame().run(new(CsDiamondRoom))
	
}