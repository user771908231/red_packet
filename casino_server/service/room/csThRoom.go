package room

var ChampionshipRoom CSThGameRoom	//锦标赛的房间

func init() {
	ChampionshipRoom.OnInit()	//初始化,开始运行
	ChampionshipRoom.Run()
}

//锦标赛
type CSThGameRoom struct {
	*ThGameRoom
	//锦标赛房间的专有属性

}

//run游戏房间
func (r *CSThGameRoom) Run() {

}
