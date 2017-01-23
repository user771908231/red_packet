package internal

import (
	"fmt"
)

func init() {
	skeleton.RegisterCommand("echo", "echo user inputs", commandEcho)
	skeleton.RegisterCommand("testmj", "测试麻将", cmdTestMj)
}

func commandEcho(args []interface{}) interface{} {
	return fmt.Sprintf("%v", args)
}

func cmdTestMj(args []interface{}) interface{} {
	cmd := args[0].(string)
	switch cmd {
	case "create":
		handlerCreateDesk(nil)
	}
	return ""
}
