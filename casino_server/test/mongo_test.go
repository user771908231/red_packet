package mongodb

import (
	"testing"
	"github.com/name5566/leaf/db/mongodb"
	"fmt"
	"casino_server/conf/casinoConf"
	"gopkg.in/mgo.v2/bson"
	"casino_server/mode"
)


func TestM(t *testing.T){
	//_TestSave(t)
	//saveWithSub(t)
	//update(t)
	_select(t)
	//saveSub2(t)
	//selectSub2(t)
	//selectSub3(t)
	//nestSeq(t)
}



func _TestSave(t *testing.T){
	t.Log("开始测试保存到数据库\n")

	// 获取连接 connection
	c, err := mongodb.Dial("localhost", 51668)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	// 获取回话 session
	s := c.Ref()
	defer c.UnRef(s)

	user := mode.User{}
	id,_ :=  c.NextSeq("test", "t_user", "id")
	user.Id = uint32(id)
	user.Mid = bson.NewObjectId();
	e := s.DB("test").C("t_user").Insert(user)
	if e != nil {
		t.Error(e)
	}

	var testResult mode.T_test
	testResult.Name = "测试2"
	testResult.Mid = bson.NewObjectId()
	s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_TEST).Insert(testResult)

	t.Log("开始测试保存到数据库--end\n")

}

func _Del(t *testing.T){
	t.Log("开始测试删除\n")
}

func update(t *testing.T){
	t.Log("开始测试保存到数据库\n")


	t.Log("\n开始测试保存到数据库--end\n")


}

func _select(t *testing.T){
	t.Log("开始测试查询一条数据\n")

	// 获取连接 connection
	c, err := mongodb.Dial("localhost", 51668)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	// 获取回话 session
	s := c.Ref()
	defer c.UnRef(s)

	result := mode.User{}
	s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_USER).Find(bson.M{"id": 19}).One(&result)
	t.Log("Mid ",result.Mid)
	t.Log("Id ",result.Id)
	t.Log("Mobile ",result.Mobile)
	t.Log("Name ",result.Name)
	//t.Log("NickName ",result.NickName)

	//result := Person{}
	//err = collection.Find(bson.M{"phone": "456"}).One(&result)
	//fmt.Println("Phone:", result.NAME, result.PHONE)

	t.Log("\n开始测试查询一条数据--end\n")


}

func _Find(t *testing.T){
	t.Log("开始测试查询多条数据\n")
}

func saveWithSub(t *testing.T){
	t.Log("开始测试保存携带子节点数据\n")
	test := mode.T_test{}
	sub := mode.T_test_sub{}
	test.Sub = sub
	test.Name = "test1"
	sub.Sname = "sname1"

	//连接数据库
	c,err := mongodb.Dial(casinoConf.DB_IP,casinoConf.DB_PORT)
	if err != nil{
		t.Error(err)
	}
	defer  c.Close()

	//获取session
	s := c.Ref()
	defer s.Close()

	s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_TEST).Insert(test)

	t.Log("\n开始测试保存携带子节点数据--end\n")
}


func saveSub2(t *testing.T){
	sub2 := &mode.T_test_sub2{
		Id:111,
		Sname:"sub2",
	}

	//连接数据库
	c,err := mongodb.Dial(casinoConf.DB_IP,casinoConf.DB_PORT)
	if err != nil{
		t.Error(err)
	}
	defer  c.Close()

	//获取session
	s := c.Ref()
	defer s.Close()
	s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_SUB2).Insert(sub2)

}
//func get

func  selectSub2(t *testing.T){
	//连接数据库
	c,err := mongodb.Dial(casinoConf.DB_IP,casinoConf.DB_PORT)
	if err != nil{
		t.Error(err)
	}
	defer  c.Close()

	//获取session
	s := c.Ref()
	defer s.Close()

	var result mode.T_test_sub2
	s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_SUB2).Find(bson.M{"sname" : "sub2"}).One(&result)
	t.Log("ObjId",result.ObjId)
	t.Log("id",result.Id)

	var testResult mode.T_test
	s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_TEST).Find(bson.M{"_id":bson.ObjectIdHex("576148b21ba69d3aa3ec410d")}).One(&testResult)
	t.Log("testResult.id",testResult.Mid)

	 e := s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_TEST).Update(bson.M{"_id":testResult.Mid}, bson.M{"$push": bson.M{ "sub2" : result.ObjId}})
	if e != nil {
		t.Error(e)
	}

}


func selectSub3(t *testing.T){
	//连接数据库
	c,err := mongodb.Dial(casinoConf.DB_IP,casinoConf.DB_PORT)
	if err != nil{
		t.Error(err)
	}
	defer  c.Close()

	//获取session
	s := c.Ref()
	defer s.Close()

	var testResult mode.T_test
	s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_TEST).Find(bson.M{"_id":bson.ObjectIdHex("576148b21ba69d3aa3ec410d")}).One(&testResult)
	t.Log("testResult.id",testResult.Mid)
	t.Log("testResult.Sub2[0]",testResult.Sub2[0])


	sub2 := &mode.T_test_sub2{}
	e := s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_SUB2).FindId(testResult.Sub2[0]).One(&sub2)
	if e != nil {
		t.Error("出错了")
		t.Error(e)
	}

	t.Log("sub2.Sname")
	t.Log("sub2.Sname",sub2.Sname)
}



func  nestSeq(t *testing.T){
	//连接数据库
	c,err := mongodb.Dial(casinoConf.DB_IP,casinoConf.DB_PORT)
	if err != nil{
		t.Error(err)
	}
	defer  c.Close()

	//获取session
	s := c.Ref()
	defer s.Close()
	id,_ :=  c.NextSeq("test", "t_user", "id")
	fmt.Println("id",id)
}




