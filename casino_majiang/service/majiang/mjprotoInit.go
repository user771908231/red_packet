package majiang

func NewMjpai() *MJPai {
	ret := &MJPai{}
	ret.Value = new(int32)
	ret.Des = new(string)
	ret.Flower = new(int32)
	ret.Index = new(int32)
	return ret

}