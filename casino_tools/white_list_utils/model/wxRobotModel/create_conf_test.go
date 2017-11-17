package wxRobotModel

import "testing"

func TestCreateConfig_GetKeywords(t *testing.T) {
	t.Log(ZzhzConf.GetKeywords("转转 翻倍 16局","红中 8局 抓 6鸟 翻倍"))
	t.Log(PdkConf.GetKeywords("十五张 20局 2人 不出黑桃3","首出黑桃3 16张 三人 10局"))

}
