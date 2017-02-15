package main

import (
	"fmt"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Language struct {
	Name string
	Language []struct{
		Name string
		Type int
		Level int
	}
}

func main() {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("runoob")

	result := Language{}
	err = c.Find(
		bson.M{"language.name": "java"},
	).Select(
		bson.M{"language.$":1},
	).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result)
}