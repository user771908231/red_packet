package main

import (
	//"github.com/garyburd/redigo/redis"
	//"fmt"
	"io/ioutil"
	"net/http"
	"math"
	"log"
	"encoding/json"
)

func main() {
	type Person struct {
		Result struct{
			       Formatted_address string
		       }
	}
	p := &Person{}
	//经度
	lng := "28.696117043877"
	//纬度
	lat := "115.95845796638"
	r, e := http.Get("http://api.map.baidu.com/geocoder/v2/?ak=DD279b2a90afdf0ae7a3796787a0742e&location="+lng+","+lat+"&output=json&pois=0")
	body, _ := ioutil.ReadAll(r.Body)
	log.Printf("r : %+v ,e : %v", string(body), e)
	err := json.Unmarshal([]byte(body), &p)
	if err != nil {
		log.Printf("解析json的时候err:%v", err)
	}
	log.Printf("得到的ret :%v", p)


	//fmt.Println(EarthDistance(39.9,116.3,31.47,104.73))

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


