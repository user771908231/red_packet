package mongo

import (
	"testing"
	"encoding/json"
)

func TestBson(t *testing.T) {
	json_str := `{
	"id":12,
	"friends":["james",
		"jack"]
		}`
	type People struct {
		Id *uint32
		Friends []*string
	}
	bson_map := People{}
	err := json.Unmarshal([]byte(json_str), &bson_map)
	t.Log(err)
	t.Log(*bson_map.Id, *bson_map.Friends[0])

	str, err := json.MarshalIndent(bson_map, "", "")
	t.Log(err)
	t.Log(string(str))

	strusts := People{}
	err = json.Unmarshal(str, &strusts)
	t.Log(err)
	t.Log(*strusts.Id, *strusts.Friends[0])
}
