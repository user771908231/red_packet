package mongodb

import (
	"testing"
	"fmt"
	"time"
	"casino_server/utils/time"
)

var a int32 = 20

func TestTemp(t *testing.T) {

	timeNow := time.Now()
	fmt.Println("time:",timeNow.Format(timeUtils.TIME_LAYOUT))
}

