package mode

import "gopkg.in/mgo.v2/bson"

type T_test struct {
	Id bson.ObjectId `json:"id" bson:"_id`
	number  uint32
	Name string
	Sub T_test_sub
	Sub2 []bson.ObjectId

}

type T_test_sub struct{
	Id uint32
	Sname string

}



type T_test_sub2 struct{
	ObjId bson.ObjectId `json:"ObjId" bson:"_id"`
	Id uint32
	Sname string

}
