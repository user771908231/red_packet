package numUtils

import (
	"fmt"
	"strconv"
)


/**
	int类型转字符串
 */
func Int2String(i int32) (string,error){
	str := fmt.Sprintf("%d", i)
	return str,nil
}

/**
 字符转转int类型
 */
func String2Int(s string) (int,error){
	i, err := strconv.Atoi(s)
	return i,err
}

