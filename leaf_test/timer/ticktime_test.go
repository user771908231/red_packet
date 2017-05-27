package timer

import (
	"time"
	"testing"
	"github.com/name5566/leaf/timer"
	"sync"
	"fmt"
)

func TestTimer(t *testing.T) {
	new_timer := timer.NewDispatcher(10)
	t.Log("1:",time.Now())
	new_timer.AfterFunc(1, func() {
		t.Log("2:",time.Now())
	})
	new_timer.AfterFunc(1, func() {
		t.Log("2:",time.Now())
	})
	(<-new_timer.ChanTimer).Cb()
	(<-new_timer.ChanTimer).Cb()
}

//模拟Cron定时任务
func TestCronTab(t *testing.T) {
	go func() {
		d := timer.NewDispatcher(10)

		// cron expr
		cronExpr, err := timer.NewCronExpr("0 0 0 * * *")
		if err != nil {
			return
		}

		d.CronFunc(cronExpr, func() {
			fmt.Println("My name is Leaf")
		})

		// dispatch
		for {
			fmt.Println(cronExpr.Next(time.Now()))
			(<-d.ChanTimer).Cb()
		}
	}()

	// Output:
	// My name is Leaf
}

//定时器
func TestInerval(t *testing.T) {
	Interval(1,10, func(second int64) bool {
		t.Log("已执行", second,"秒")
		if second >= 5 {
			return false
		}
		return true
	})
}

//定时器
func Interval(step int64, length int64, fn func(int64)bool){
	start_time := time.Now().Unix()
	var after_func func()
	var wg sync.WaitGroup
	after_func = func() {
		now := time.Now().Unix()
		if now < start_time + length {
			second := now-start_time
			is_continue := fn(second)
			if is_continue {
				wg.Add(1)
				time.AfterFunc(time.Duration(step*int64(time.Second)), after_func)
			}
		}
		wg.Done()
	}
	wg.Add(1)
	time.AfterFunc(time.Duration(step*int64(time.Second)), after_func)
	wg.Wait()
}
