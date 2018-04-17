package adminModel

import (
	"gopkg.in/mgo.v2/bson"

	"time"
	"crypto/md5"
	"encoding/hex"
	"casino_common/utils/db"
	"sendlinks/conf/tableName"
)
//后台用户
type Admin struct {
	ObjIds					bson.ObjectId	`bson:"_id"`
	NickName				string	//昵称
	AccountName				string	//帐户
	Password				string	//密码
	AccountPower			[]*Power	//权限列
	Time 					time.Time
}

func (A *Admin)  Insert() error {
	A.ObjIds = bson.NewObjectId()
	p := md5.New()
	p.Write([]byte(A.Password))
	A.Password = hex.EncodeToString(p.Sum(nil))
	A.Time = time.Now()
	err :=  db.C(tableName.DB_LINKS_ADMIN_INFO).Insert(A)
	return err
}
//修改
func (A *Admin) Save() error {
	err := db.C(tableName.DB_LINKS_ADMIN_INFO).Update(bson.M{"_id": A.ObjIds}, A)
	return err
}
//删除
func DelAdminUser(s string) error {
	err := db.C(tableName.DB_LINKS_ADMIN_INFO).Remove(bson.M{"_id":bson.ObjectIdHex(s)})
	return err
}
//查询
