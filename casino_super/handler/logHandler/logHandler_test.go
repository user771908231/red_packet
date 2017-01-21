package logHandler

import (
	"testing"
	"fmt"
)

func TestPaginator(t *testing.T) {
	maps := Paginator(10, 100, 3900)

	println(fmt.Sprintf("pages %v", maps["pages"]))
	println(fmt.Sprintf("totalpages %v", maps["totalpages"]))
	println(fmt.Sprintf("firstpage %v", maps["firstpage"]))
	println(fmt.Sprintf("lastpage %v", maps["lastpage"]))
	println(fmt.Sprintf("currpage %v", maps["currpage"]))
}

