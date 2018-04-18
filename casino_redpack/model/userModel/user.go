package userModel

import (
	"casino_redpack/conf/config"
	"casino_common/utils/db"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"fmt"
	time "time"
	"casino_common/common/userService"
	"casino_common/common/log"
)

const USER_TABLE_NAME string = "t_redpack_user"

type User struct {
	Id         uint32
	ThreePartyId uint32 //平台注册的ID
	Level      int32
	NickName   string
	HeadUrl    string
	OpenId     string
	UnionId    string
	PassWd     string
	SignUpTime time.Time
	Coin       float64 //金币
	ExtensionId	int32
}

func (User *User) CapitalUplete(action string ,munber float64,msg string) error{
	if action == "-" {
		if User.ThreePartyId != 0 {
			log.T("平台用户金币减")
			_,errr:= userService.DECRUserCOIN(User.ThreePartyId,int64(munber),msg)
			if errr != nil {
				return errr
			}
		}

		err := User.findAndModify(bson.M{"coin": -munber})
		if err != nil{
			return err
		}
		return nil
	}else if action == "+" {
		if User.ThreePartyId != 0 {
			log.T("平台用户金币加")
			_,errr:= userService.INCRUserCOIN(User.ThreePartyId,int64(munber),msg)
			if errr != nil {
				return errr
			}
		}
		err := User.findAndModify(bson.M{"coin": munber})
		if err != nil{
			return err
		}
		return nil
	}else if action == "/" {
		User.Coin = User.Coin / munber
		err :=  NextReturn(User)
		if err != nil{
			return err
		}
		return nil
	}else if action == "*" {
		User.Coin = User.Coin * munber
		err :=  NextReturn(User)
		if err != nil{
			return err
		}
		return nil
	}else {
		return nil
	}
}

func(U *User) findAndModify(Update interface{}) error {
	err := db.C(USER_TABLE_NAME).Update(bson.M{"id":U.Id},bson.M{"$inc": Update } )
	if err != nil {
		return err
	}
	return nil
}

func (U *User) Uplate() error {
	err := db.C(USER_TABLE_NAME).Update(bson.M{"id":U.Id},U)
	if err != nil {
		return err
	}
	return nil
}

func NextReturn(User *User)error  {
	err := db.C(USER_TABLE_NAME).Update(bson.M{"id":User.Id},User)
	if err != nil {
		return err
	}
	return nil
}

//通过ThreePartyId获取用户信息
func GetUserByThreePartyId(id uint32) *User {
	var err error = nil
	user_row := new(User)
	err = db.C(USER_TABLE_NAME).Find(bson.M{
		"threepartyid": id,
	}, user_row)
	if err != nil {
		return nil
	}
	return user_row
}

//通过id获取用户信息
func GetUserById(id uint32) *User {
	var err error = nil
	user_row := new(User)
	err = db.C(USER_TABLE_NAME).Find(bson.M{
		"id": id,
	}, user_row)
	if err != nil {
		return nil
	}
	return user_row
}

//验证密码密码
func Login(user_name string, passwd string) *User {
	var err error = nil
	user_row := new(User)
	h := md5.New()
	h.Write([]byte(passwd))
	passwd = hex.EncodeToString(h.Sum(nil))
	err = db.C(USER_TABLE_NAME).Find(bson.M{
		"nickname": user_name,
		"passwd":   passwd,
	}, user_row)
	if err != nil {
		return nil
	}
	return user_row
}

//验证手机号 密码
func   TableValues(user_name string , passwd_one string , passwd_two string,extension_id int32) (error,string){
	_,Msg := JudgeMobilePhoneWhetherSignuo(user_name)
	if Msg != "" {
		return  nil,"该手机号已注册"
	}
	_,errorMsg := VerificationRegexp(user_name)
	if errorMsg == ""  &&  passwd_one == passwd_two {

		dData := User{
			NickName:user_name,
			PassWd:passwd_one,
			ExtensionId:extension_id,
		}
		err := dData.Insert()
		if err != nil {
			return err,"新建用户失败"
		}
		return nil,""
	}

	return nil,"账号或密码验证错误！"

}
//判断手机是否已注册

func JudgeMobilePhoneWhetherSignuo(MobilePone string)  (error,string){
	var err error = nil
	MobilePoneRow := new(User)
	err = db.C(USER_TABLE_NAME).Find(bson.M{"nickname": MobilePone}, MobilePoneRow)
	if err != nil {
		return nil,""
	}
	if MobilePoneRow != nil {
		return nil,"该手机号已注册"
	}
	return err,""
}

//正则
func VerificationRegexp(str string) (result string, errorMsg string){
	//移动：139   138   137   136   135   134   147   150   151   152   157   158    159   178  182   183   184   187   188
	//联通：130   131   132   155   156   185   186   145   176
	//电信：133   153   177   173   180   181   189
	//虚拟运营商：170  171
	r,_:=regexp.Compile("/^0?(13[0-9]|15[012356789]|17[013678]|18[0-9]|14[57])[0-9]{8}$/")
	b := r.MatchString(str)
	if b != false {
		dData := VerificationErrorMessage{
			msg:str,
		}
		errorMsg = dData.Error()
		return
	}else {
		return str,""
	}
}

type VerificationErrorMessage struct { msg string }

func (vem *VerificationErrorMessage) Error() string {
	ErrMessage := `
		电话格式不对
		msg : %d
	`
	return fmt.Sprintf(ErrMessage,vem.msg)
}


//插入一个新用户
func (user *User) Insert() error{
	id, err := db.GetNextSeq(config.DB_USER_SEQ)
	if err != nil {
		return errors.New("获取user_id自增键失败！")
	}
	user.Id = uint32(id)
	h := md5.New()
	h.Write([]byte(user.PassWd))
	user.PassWd = hex.EncodeToString(h.Sum(nil))
	user.SignUpTime = time.Now()
	err = db.C(USER_TABLE_NAME).Insert(user)
	return err
}

//编辑用户资料
func (user *User) Save() error {
	err := db.C(USER_TABLE_NAME).Update(bson.M{"id": user.Id}, user)
	return err
}
//获取全部用户
func GetUsers(query bson.M) []*User {
	users := []*User{}
	err := db.C(USER_TABLE_NAME).FindAll(query,&users)
	if err != nil {
		return nil
	}
	return users
}


