package keysModel

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"casino_common/utils/db"
	"sendlinks/conf/tableName"
)

type Keys struct {
	ObjId 	bson.ObjectId	`bson:"_id"`
	Keys	string
	Remarks	string
	Status	int	//0 关 1 开
	Time 	time.Time
	KeysNumber bson.ObjectId
}

//分表记录
type KeysNumber struct {
	ObjId bson.ObjectId `bson:"_id"`
	Number uint32
	Time time.Time
	KeyListItem	[]*KeyListItem
}

func (K *KeysNumber)	Insert() error {
	K.ObjId = bson.NewObjectId()
	K.Time = time.Now()
	err := db.C(tableName.DB_KEYS_NUMBER).Insert(K)
	return err

}
//保持KeysNumber与Keys数据一致性
type KeyListItem struct {
	ObjId bson.ObjectId `bson:"_id"`
	FirstTime	time.Time
	UpdateTime	time.Time
}

func (K *Keys)  Insert() error {
	K.ObjId = bson.NewObjectId()
	K.Status = 1
	K.Time = time.Now()
	K.KeysNumber = IsTablesNumber()
	err := db.C(tableName.DB_KEYS_LISTS).Insert(K)
	return err
}

func IsTablesNumber() bson.ObjectId {
	lits :=  new([]*KeysNumber)
	db.C(tableName.DB_KEYS_NUMBER).FindAll("",lits)
	//判断错误 和lits为nil
	if lits == nil {
		//创建一条数据
		if err,ObjId := CreateKeysNumber();err == nil {
			return ObjId
		}
		_,Id := CreateKeysNumber()
		return Id
	}else {
		ObjId := IsTableNumber(*lits)
			return ObjId
		
	}

}

func CreateKeysNumber() (error,bson.ObjectId){
	k := new(*KeysNumber)
	err := db.C(tableName.DB_KEYS_NUMBER).Insert(k)
	return err, k.ObjId
}
//获取TableNumber不足100000条的表 
func IsTableNumber(K []*KeysNumber) bson.ObjectId {
	lengt := len(K)
	for i,itme := range K{
		if i == lengt {
			continue
		}
		if itme.Number > uint32(100000) {
			return itme.ObjId
		}
	}
	_,Id := CreateKeysNumber()
	return Id
}

func GetKeysNumber(Id bson.ObjectId) *KeysNumber {
	v := new( KeysNumber)
	err := db.C(tableName.DB_KEYS_NUMBER).Find(bson.M{"_id":Id},v)
	if err != nil {
		return nil
	}
	return v
}

func GetKeys() []*Keys {
	KN := []*KeysNumber{}
	K := []*Keys{}
	db.C(tableName.DB_KEYS_NUMBER).FindAll(bson.M{},KN)
	if KN != nil {
		lengt := len(KN)
		if lengt != 0 {
			for i,item := range KN {
				if i == lengt {
					continue
				}
				db.C(tableName.DB_KEYS_LISTS).FindAll(bson.M{},K)
				K = append(K,K)
			}
		}

	}
	return K
}
