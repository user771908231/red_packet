package mode

import "gopkg.in/mgo.v2/bson"

type T_test struct {
	Mid		bson.ObjectId		`json:"mid" bson:"_id"`
	id uint32
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
	ObjId bson.ObjectId 	`json:"ObjId" bson:"_id"`
	Id uint32
	Sname string

}
