package room

type CSDiamondGame struct {
	*CsGameSkeleton
}

func NewCsDiamondGame() CSGame {
	game := &CSDiamondGame{}
	return game
}