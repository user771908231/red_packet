package doudizhu

//备份用
func NewPDdzbak() *PDdzbak {
	ret := new(PDdzbak)
	return ret
}

//Desk
func NewPDdzDesk() *PDdzDesk {
	ret := new(PDdzDesk)
	return ret
}

func NewPDdzUser() *PDdzUser {
	ret := new(PDdzUser)
	return ret
}


//New一个Desk
func NewDdzDesk() *DdzDesk {
	desk := new(DdzDesk)
	desk.PDdzDesk = NewPDdzDesk()
	return desk
}

func NewDdzUser() *DdzUser {
	user := new(DdzUser)
	user.PDdzUser = NewPDdzUser()
	return user
}

