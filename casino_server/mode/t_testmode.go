package mode


type T_test struct {
	Id uint32
	Name string
	Sub T_test_sub
}

type T_test_sub struct{
	Id uint32
	Sname string

}

