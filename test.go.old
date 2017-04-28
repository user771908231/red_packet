package main

import (
	"io/ioutil"
	"net/http"
	"log"
	"encoding/json"

)

func main() {
	type Person struct {
		Code int
		Data    struct{
			Country string
			Area_id string
			Isp_id string
			Ip string
			Region string
			City string
			Isp string
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

}
