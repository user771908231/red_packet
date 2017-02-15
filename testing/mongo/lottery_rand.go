package main

import (
	"math/rand"
	"log"
	"time"
)

func main() {
	//arr_goods := []string{"10%","20%","30%","40%","50%"}
	arr_chance := []int{10,20,30,40,50}

	//统计抽奖结果
	test_res := make(map[int]int)
	for i:=0;i<10000;i++ {
		rand_index := GetRandIndex(arr_chance)
		rand_chance := arr_chance[rand_index]
		if _,ok := test_res[rand_chance];ok {
			test_res[rand_chance] += 1
		}else {
			test_res[rand_chance] = 1
		}
	}

	log.Println(test_res)

}

//得到概率数组的索引
func GetRandIndex(arr []int) int {
	count := 0
	for _,num := range arr {
		count += num
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand_num := r.Intn(count)

	min,max := 0,0
	for i,chance := range arr {
		max = min+chance
		if rand_num >= min && rand_num < max {
			return i
		}
		min += chance
	}

	return 0
}
