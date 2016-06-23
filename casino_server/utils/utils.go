package utils

import (
	"math/rand"
	"os"
	"strconv"
	"time"
	"casino_server/common/log"
	"github.com/golang/protobuf/proto"
)

func WriteFile(data []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Error("Failed to create the output file: %v", filename)
		return err
	}
	defer file.Close()
	_, err = file.Write(data)

	if err == nil {
		log.T("write to file %v ret ok", filename)
	} else {
		log.Error("write to file %v failed, err:%v", filename, err)
	}
	//for _, value := range values {
	//	str := strconv.Itoa(value)
	//	file.WriteString(str)
	//}
	return err
}

func ReadFile(filename string) (data []byte, err error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Error("Failed to open file: %v", filename)
		return nil, err
	}
	defer file.Close()
	data = make([]byte, 1024000)
	n, err := file.Read(data)
	log.T("file.Read ret:%v err:%v", n, err)
	return data[:n], err
}

func Itoa(n int) string {
	return strconv.FormatInt(int64(n), 10)
}

func Ntoa(n int32) string {
	return strconv.FormatInt(int64(n), 10)
}

func Utoa(n uint32) string {
	return strconv.FormatUint(uint64(n), 10)
}

func Atou(s string) uint32 {
	u64, _ := strconv.ParseUint(s, 10, 0)
	return uint32(u64)
}

func Atoi(s string) int32 {
	i, _ := strconv.Atoi(s)
	return int32(i)
}

func Atof(s string) float32 {
	f, _ := strconv.ParseFloat(s, 32)
	return float32(f)
}

func Now() uint32 {
	return uint32(time.Now().UTC().Unix())
}

func TodayDawn() uint32 {

	year, month, day := time.Now().UTC().Date()

	timeValue := string(strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day) + " 00:00:00")
	if getTime, err := time.Parse("2006-1-2 15:04:05", timeValue); err == nil {
		return uint32(getTime.Unix())
	}

	return 0
}

//return whether t is today
func IsToday(t uint32) bool {
	y1, m1, d1 := time.Unix(int64(t), 0).UTC().Date()
	y2, m2, d2 := time.Now().UTC().Date()

	log.T(" IsToday() ->  lastTime: %v-%v-%v   today: %v-%v-%v (%v)", y1, m1, d1, y2, m2, d2, t)
	if y1 == y2 && m1 == m2 && d1 == d2 { //is same day
		return true
	}

	return false
}

func IsYestoday(t uint32) bool {

	return IsToday(t + 86400)
}

func WeekDay() time.Weekday {
	return time.Now().UTC().Weekday()
}

//return rand number between [start, end)
func Rand(start, end int32) int32 {
	if end-start <= 0 {
		return start
	}
	return start + rand.Int31n(end-start)
}

//return rand number between [0, n)
func Randn(n int32) int32 {
	if n <= 0 {
		return 0
	}
	return rand.Int31n(n)
}

//rate range from: [0 - 1.0]
func HitRandomRate(rate float32) bool {
	randnum := Randn(100)
	log.T("HitRandomRate :: [ randnum %v ? rate: %v ] ==> %v", randnum, rate, randnum < int32(rate*100))
	return randnum < int32(rate*100)
}

func RemoveSlice(s []interface{}, pos int) {
	s = append(s[:pos], s[pos+1:]...)
}

func RemovePSlice(s []proto.Message, pos int) {
	//	s = append(s[:pos], s[pos+1:]...)
}
