package test

import (
	"testing"
	"fmt"
	"casino_server/utils/jobUtils"
	"time"
)

func TestMongoUtils(t *testing.T) {
	jobUtils.DoAsynJob(time.Second * 1, func() bool {
		fmt.Println("测试定时任务...")
		return false
	})

	for ; ;  {
		
	}
}
