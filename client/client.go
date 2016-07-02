package main

import (
)
import "client/a"

const url  = "192.168.199.120:3563"
const TCP = "tcp"

func main() {
	a := &a.A{}
	a.Aa()
	a.Ba()
}