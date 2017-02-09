package huParserIns

type ChengDuHuParser struct {
	*HuParserSkeleton
}

func NewChengDuHuParser() *ChengDuHuParser {
	p := &ChengDuHuParser{
		HuParserSkeleton: NewHuParserSkeleton(),
	}
	return p
}
