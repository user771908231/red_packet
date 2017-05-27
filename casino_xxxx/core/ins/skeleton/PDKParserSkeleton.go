package skeleton

type PDKParserSkeleton struct {
}

//todo 出牌的时候比较大小
func (p *PDKParserSkeleton) CanOut(outpai interface{}, check interface{}) (bool, error) {
	return false, nil
}

//todo 通过一副牌的id解析牌型
func (p *PDKParserSkeleton) Parse(pids []int32) (interface{}, error) {
	return false, nil
}
