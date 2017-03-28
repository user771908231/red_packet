package dao

import "testing"

func TestGetMJRoomConfigData(t *testing.T) {
	data := GetMJRoomConfigData()
	t.Log(data)
	for _, d := range data {
		t.Log(d)
	}
}
