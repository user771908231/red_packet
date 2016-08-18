package bbproto

//返回一个userSeesion
func NewThServerUserSession() *ThServerUserSession {
	result := &ThServerUserSession{}
	result.DeskId = new(int32)
	result.MatchId = new(int32)
	result.UserId = new(uint32)
	result.GameType = new(int32)
	result.GameStatus = new(int32)
	return result
}