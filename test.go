package main

import (
	//"github.com/garyburd/redigo/redis"
	//"fmt"
	"io/ioutil"
	"net/http"
	"math"
	"log"
	"encoding/json"
	"fmt"
)

func main() {
	type Person struct {
		Result struct{
			       Formatted_address string
		       }
	}
	p := &Person{}
	lat:= "30.548635"
	lng:= "104.06328"
	r, e := http.Get("http://api.map.baidu.com/geocoder/v2/?ak=IBeISijoITQMKkuK9Dlp56ELYP2a8Nek&location="+lat+","+lng+"&output=json&pois=0")
	body, _ := ioutil.ReadAll(r.Body)
	log.Printf("r : %+v ,e : %v", string(body), e)
	err := json.Unmarshal([]byte(body), &p)
	if err != nil {
		log.Printf("解析json的时候err:%v", err)
	}
	log.Printf("得到的ret :%v", p)


	//fmt.Println(EarthDistance(30.548646926879883,104.06388092041016,30.548696517944336,104.0637435913086))

	for i:=0;i<5;i++{
		fmt.Println(i)
	}

	m := map[string]string{
		"a" : "aa" ,
		"b" : "bb" ,
	}

	for k,v := range m{
		fmt.Println(k,v)
	}
}

//经纬度计算距离
func EarthDistance(lat1, lng1, lat2, lng2 float64) float64 {
	radius := 6378137.00 // 6371000
	rad := math.Pi/180.0

	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad

	theta := lng2 - lng1
	dist :=math.Acos(math.Sin(lat1) * math.Sin(lat2) + math.Cos(lat1) * math.Cos(lat2) * math.Cos(theta))
	err := radius * dist
	return err
}


