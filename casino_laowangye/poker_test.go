package main

import (
	"math/rand"
	"time"
	"testing"
	"fmt"
	"casino_common/utils/numUtils"
	"casino_laowangye/service/laowangye"
	"casino_common/proto/ddproto"
)

//色子类型解析概率统计
func TestShaiziTypeParseCount(t *testing.T) {
	sum := 0
	num_map := map[ddproto.LwyShaiziType]int32{
		ddproto.LwyShaiziType_LWY_SHAIZI_TYPE_7_DIAN: 0,
		ddproto.LwyShaiziType_LWY_SHAIZI_TYPE_8_DIAN: 0,
		ddproto.LwyShaiziType_LWY_SHAIZI_TYPE_LAN_DIAN: 0,
	}
	for i:=1111; i<6667; i++ {
		str := []byte(fmt.Sprintf("%d", i))

		list := []int32{int32(numUtils.String2Int(string(str[0]))), int32(numUtils.String2Int(string(str[1]))), int32(numUtils.String2Int(string(str[2]))), int32(numUtils.String2Int(string(str[3])))}

		shaizi_type,err := laowangye.ParseShaiziType(list)
		if err == nil || err.Error() != "点数非法" {
			sum++
			fmt.Printf("%4d %v %v %v \n",sum, list, shaizi_type, err)
			num_map[shaizi_type]++
		}
	}
	fmt.Println(num_map)
}

//色子解析
func TestShaiziParse(t *testing.T) {
	list := make([]int32, 4)
	for i:=0; i<4; i++ {
		list[i] = rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(6) + 1
	}
	shaizi_type, err := laowangye.ParseShaiziType(list)
	fmt.Println(list, shaizi_type, err)
}
