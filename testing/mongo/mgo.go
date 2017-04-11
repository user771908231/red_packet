package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//用户
type User struct {
	Name string
}

func main() {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("users")

	result := []User{}

	err = c.Find(bson.M{
		"name": "james",
	}).All(&result)

	if err != nil {
		fmt.Println("err !", err.Error())
		return
	}

	fmt.Println("result:", result)
}
