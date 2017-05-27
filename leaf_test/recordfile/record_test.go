package recordfile

import (
	"testing"
	"github.com/name5566/leaf/recordfile"
	"log"
)

type Record struct {
	// index 0
	IndexInt int "index"
	// index 1
	IndexStr string "index"
	_Number  int32
	Str      string
	Arr1     [2]int
	Arr2     [3][2]int
	Arr3     []int
	St       struct {
			 Name string "name"
			 Num  int    "num"
		 }
	//M map[string]int
}

func TestRecord(t *testing.T) {
	rf, err := recordfile.New(Record{})
	if err != nil {
		log.Println("err:",err)
		return
	}

	err = rf.Read("./test.txt")
	if err != nil {
		log.Println("err1")
		return
	}

	for i := 0; i < rf.NumRecord(); i++ {
		r := rf.Record(i).(*Record)
		log.Println(r.IndexInt)
	}

	r := rf.Index(2).(*Record)
	log.Println(r.Str)

	r = rf.Indexes(0)[2].(*Record)
	log.Println(r.Str)

	r = rf.Indexes(1)["three"].(*Record)
	log.Println(r.Str)
	log.Println(r.Arr1[1])
	log.Println(r.Arr2[2][0])
	log.Println(r.Arr3[0])
	log.Println(r.St.Name)
	//log.Println(r.M["key6"])

	// Output:
	// 1
	// 2
	// 3
	// cat
	// cat
	// book
	// 6
	// 4
	// 6
	// name5566
	// 6
}

//测试RecordWrite
func TestRecordReader(t *testing.T) {
	rf,err := recordfile.New(Record{})
	if err != nil{
		t.Log("err:",err)
	}
	err = rf.Read("./test.txt")
	if err != nil {
		t.Log("err:",err)
	}
	for i:=0;i< rf.NumRecord();i++ {
		row := rf.Record(i).(*Record)
		t.Log(*row)
	}

}