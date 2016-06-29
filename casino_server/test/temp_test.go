package mongodb

import (
	"testing"
	"fmt"
	"time"
)

func TestTemp(t *testing.T) {
	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for _ = range ticker.C {
			fmt.Println("ticked at %v", time.Now())
		}
	}()

	for ; ;  {
		
	}
}

