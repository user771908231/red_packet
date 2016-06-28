package utils

import "fmt"


/**

 */
func Int2String(i int32) (string,error){
	str := fmt.Sprintf("%d", i)
	return str,nil
}
