package keysModel

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"casino_common/utils/db"
	"sendlinks/conf/tableName"

	"os"
	"bufio"
	"fmt"
	"errors"
	"strings"
)

type Keys struct {
	ObjId 	bson.ObjectId	`bson:"_id"`
	Keys	string
	Remarks	string
	Status	int	//0 关 1 开
	Time 	time.Time
}

func (K *Keys)  Insert() error {
	K.ObjId = bson.NewObjectId()
	K.Status = 0
	K.Time = time.Now()
	err := db.C(tableName.DB_KEYS_LISTS).Insert(K)
	return err
}

func (K *Keys)  Update() error {
	err := db.C(tableName.DB_KEYS_LISTS).Update(bson.M{"_id":K.ObjId},K)
	return err
}

func (K *Keys)  Del() error {
	err := db.C(tableName.DB_KEYS_LISTS).Remove(bson.M{"_id":K.ObjId})
	return err
}
func GetKeysAll(query bson.M,page int,number int) (int,[]*Keys){
	list := []*Keys{}
	_,count := db.C(tableName.DB_KEYS_LISTS).Page(query, &list, "-requesttime", page, number)
	return count,list
}

func IdKeyRow(string string) *Keys {
	row := new(Keys)
	err := db.C(tableName.DB_KEYS_LISTS).Find(bson.M{"_id":bson.ObjectIdHex(string)},row)
	if err != nil {
		return nil
	}
	return row
}

func KeysStatus(string string,status int) error {
	row := GetkeysId(string)
	row.Status = status
	err :=row.Update()
	return err
}

func GetkeysId(string string) *Keys {
	row := new(Keys)
	err := db.C(tableName.DB_KEYS_LISTS).Find(bson.M{"_id":bson.ObjectIdHex(string)},row)
	if err != nil {
		return nil
	}
	return row
}


func (K *Keys)  Upsert() error {
	err := db.C(tableName.DB_KEYS_LISTS).Upsert(bson.M{"keys":K.Keys},K)
	return err
}

//读取文件
func OpenFiles(string string){
	f,err := os.Open(string)
	if err != nil {
		errors.New("打开文件错误！")
	}
	defer f.Close()
	b := bufio.NewReader(f)
	line, err := b.ReadString(',')
	for ; err == nil; line, err = b.ReadString(',') {
		K := Keys{
			Keys:strings.Trim(line,","),
		}
		K.Insert()
		//if err != nil {
		//	log.E("失败！%s",err)
		//}else{
		//	log.E("成功！")
		//}
		//arr = append(arr,str)
	}
	del := os.Remove(string)
	if del != nil {
		fmt.Println(del);
	}

	//if err == io.EOF {
	//	fmt.Print(line)
	//} else {
	//	return errors.New("read occur error!")
	//}
	//return nil
}



