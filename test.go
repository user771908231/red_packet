package main

import (
	"io/ioutil"
	"net/http"
	"log"
	"encoding/json"
	//"math"
	//"fmt"
)

func main() {
	type Person struct {
		Data    struct{
			City string
			}
	}
	//unmarshal to struct
	ip :="180.89.94.90"
	p := &Person{}
	r, e := http.Get("http://ip.taobao.com/service/getIpInfo.php?ip="+ip)
	body, _ := ioutil.ReadAll(r.Body)
	log.Printf("r : %+v ,e : %v", string(body), e)
	err := json.Unmarshal([]byte(body), &p)
	if err != nil {
		log.Printf("解析json的时候err:%v", err)
	}
	log.Printf("得到的ret :%v", p)




	//lat1 := 29.490295
	//lng1 := 106.486654
	//
	//lat2 := 29.615467
	//lng2 := 106.581515
	//fmt.Println(EarthDistance(lat1, lng1, lat2, lng2))

}

//func EarthDistance(lat1, lng1, lat2, lng2 float64) float64 {
//	radius := 6371000.00 // 6378137
//	rad := math.Pi/180.0
//
//	lat1 = lat1 * rad
//	lng1 = lng1 * rad
//	lat2 = lat2 * rad
//	lng2 = lng2 * rad
//
//	theta := lng2 - lng1
//	dist :=math.Acos(math.Sin(lat1) * math.Sin(lat2) + math.Cos(lat1) * math.Cos(lat2) * math.Cos(theta))
//	err := radius * dist
//	return err
//}
