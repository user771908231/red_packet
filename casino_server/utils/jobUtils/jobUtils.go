package jobUtils

import "time"

//定义一个异步任务
func DoAsynJob(d time.Duration, f func() bool) {
	ticker := time.NewTicker(d)
	go func() {
		for _ = range ticker.C {
			if f() {
				break;
			}
		}
	}()
}