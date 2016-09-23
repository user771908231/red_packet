package room

type CSDiamondGame struct {
	*CsGameSkeleton
}

func NewCsDiamondGame() *CSDiamondGame {
	game := &CSDiamondGame{}
	game.CsGameSkeleton = new(CsGameSkeleton)
	return game
}