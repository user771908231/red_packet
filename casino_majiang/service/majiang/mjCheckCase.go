package majiang

var CHECK_CASE_STATUS_0 int32 = 0        //表示没有判定过
var CHECK_CASE_STATUS_1 int32 = 1        //表示胡牌判定过
var CHECK_CASE_STATUS_2 int32 = 2        //表示碰／杠 判定过


func (c *CheckBean) IsChecked() bool {
	return false;
}