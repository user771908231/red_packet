package redpack

import (
	"casino_redpack/modules"
	"gopkg.in/mgo.v2/bson"
	"time"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
)

//绑定银行卡
func BindBankHandler(ctx *modules.Context) {

	ctx.HTML(200, "redpack/home/bind_bank")
}

//银行卡列表
func BankListHandler(ctx *modules.Context) {

	data := GetUserBanck(ctx.IsLogin().Id)

	ctx.JSON(200,bson.M{"code": 1,
		"message": "success",
		"request":  data})
}

//银行卡添加
func BankAddHandler(ctx *modules.Context) {
	ctx.HTML(200, "redpack/home/bank_add")
}

type UserBanck struct {
	Id bson.ObjectId `bson:"_id"`
	UserId uint32
	AcctNo string
	AcctName string
	AcctBankName string
	RecBankName string
	Time time.Time
}

func (U *UserBanck) Inerst() error {
	U.Id = bson.NewObjectId()
	U.Time = time.Now()
	err := db.C(tableName.TABLE_REDPACK_BANCK_LOG).Upsert(bson.M{"_id":U.Id},U)
	return err
}

func GetUserBanck(id uint32) []*UserBanck {
	data := []*UserBanck{}
	err := db.C(tableName.TABLE_REDPACK_BANCK_LOG).FindAll(bson.M{"userid":id},&data)
	if err != nil {
		return nil
	}
	return data
}

func BancklogHandler(ctx *modules.Context) {
	acct_no := ctx.Query("acct_no")
	acct_name := ctx.Query("acct_name")
	acct_bank_name := ctx.Query("acct_bank_name")
	rec_bank_name := ctx.Query("rec_bank_name")
	if acct_no == "" || acct_name == "" || acct_bank_name == "" || rec_bank_name == "" {
		ctx.Data = bson.M{
			"acct_no":acct_no,
			"acct_name" :acct_name,
			"acct_bank_name":acct_bank_name,
			"rec_bank_name" :rec_bank_name,
		}
		ctx.JSON(200,bson.M{
			"request":bson.M{
				"code":"0",
					   "msg":"填写信息不完整！",
			},

		})
		return
	}
	U := new(UserBanck)
	U.AcctNo = acct_no
	U.AcctName = acct_name
	U.AcctBankName = acct_bank_name
	U.RecBankName = rec_bank_name
	U.UserId = ctx.IsLogin().Id
	err := U.Inerst()
	if err != nil {
		ctx.JSON(200,bson.M{
			"request":bson.M{
				"code":"0",
				"msg":"新增银行卡失败！",
			},

		})
		return
	}
	ctx.JSON(200,bson.M{
		"request":bson.M{
			"code":"1",
			"msg":"填写信息不完整！",
		},

	})
	return
}
