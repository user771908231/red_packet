package test

import (
	"testing"
	"casino_doudizhu/service/doudizhu"
	"casino_server/utils/redisUtils"
	"fmt"
	"time"
	"casino_server/common/log"
)

func TestMjRedis(t *testing.T) {
	user := new(doudizhu.DdzUser)
	user.Id = new(int32)
	user.User = new(string)

	*user.Id = 329090
	*user.User = "哈哈哈"

	user.Test = new(doudizhu.TestModel)
	user.Test.TId = new(int32)
	user.Test.TUser = new(string)

	*user.Test.TId = 997678
	*user.Test.TUser = "哈哈哈"

	log.T("保存前user:[%v]", user)
	log.T("保存前user.Id:[%v]", user.Id)
	log.T("保存前user.User:[%v]", user.User)
	log.T("保存前user.Test:[%v]", user.Test)
	log.T("保存前user.Test.TId:[%v]", user.Test.TId)
	log.T("保存前user.Test.TUser:[%v]", user.Test.TUser)

	log.T("保存前user:[%v]", *user)
	log.T("保存前user.Id:[%v]", *user.Id)
	log.T("保存前user.User:[%v]", *user.User)
	log.T("保存前user.Test:[%v]", *user.Test)
	log.T("保存前user.Test.TId:[%v]", *user.Test.TId)
	log.T("保存前user.Test.TUser:[%v]", *user.Test.TUser)

	err := redisUtils.SetObj("9900", user)
	fmt.Println("查看错误", err)

	time.Sleep(time.Second * 3)

	rdata := redisUtils.GetObj("9900", &doudizhu.DdzUser{})
	if rdata == nil {
		log.E("没有从redis中找到数据")
	} else {
		ruser := rdata.(*doudizhu.DdzUser)
		log.T("保存后user:[%v]", ruser)
		log.T("保存后ruser.Id:[%v]", ruser.Id)
		log.T("保存后ruser.ruser:[%v]", ruser.User)
		log.T("保存后ruser.Test:[%v]", ruser.Test)
		log.T("保存后ruser.Test.TId:[%v]", ruser.Test.TId)
		log.T("保存后ruser.Test.Truser:[%v]", ruser.Test.TUser)

		log.T("保存后ruser:[%v]", *ruser)
		log.T("保存后ruser.Id:[%v]", *ruser.Id)
		log.T("保存后ruser.ruser:[%v]", *ruser.User)
		log.T("保存后ruser.Test:[%v]", *ruser.Test)
		log.T("保存后ruser.Test.TId:[%v]", *ruser.Test.TId)
		log.T("保存后ruser.Test.Truser:[%v]", *ruser.Test.TUser)

	}

	time.Sleep(time.Second * 10)

}