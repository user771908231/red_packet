package test

import (
	"testing"
	"casino_server/utils/pokerUtils"
	"fmt"
)

func TestPoker(t *testing.T) {
	_, rmapdes, pvalue, pflower, pname := pokerUtils.ParseByIndex(1)
	fmt.Println(rmapdes, pvalue, pflower, pname)
}