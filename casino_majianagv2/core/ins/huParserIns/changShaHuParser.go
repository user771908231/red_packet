package huParserIns

type ChangShaHuParser struct {
	*HuParserSkeleton
}

func NewChangShaHuParser() *ChangShaHuParser {
	p := &ChangShaHuParser{
		HuParserSkeleton: NewHuParserSkeleton(),
	}
	return p
}
