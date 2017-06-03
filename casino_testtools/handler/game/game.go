package game

import (
	"bufio"
	"casino_common/utils/numUtils"
	"casino_common/utils/testUtils"
	"casino_testtools/modules"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func GameTest(ctx *modules.Context) {
	sorce := ctx.Query("source")
	if sorce == ""{
		userFile := "./4/xipai.json"
		fout,err := os.Create(userFile)
		defer fout.Close()
		if err != nil {
			fmt.Println(userFile,err)
			return
		}
		fout.WriteString("")
	}
	ctx.HTML(200, "game/game")
}

func GameEdit(ctx *modules.Context) {
	gameId := ctx.Query("gameid")

	outputFile, outputError := os.OpenFile("./"+gameId+"/xipai.json",
		os.O_WRONLY|os.O_CREATE, 0666) //0666是标准的权限掩码,关于打开标识看下面
	if outputError != nil {
		fmt.Printf("An error occurred with file creation:%v\n", outputError)
		return
	}
	defer outputFile.Close()

	outputString := ctx.Query("game")
	s := strings.Split(outputString, ",")

	xipai := &testUtils.XiPai{}
	for _, s2 := range s {
		xipai.Ids = append(xipai.Ids, numUtils.String2Int(s2))
	}

	b, e := json.Marshal(xipai)
	if e != nil {
		return
	}
	outputWriter := bufio.NewWriter(outputFile)
	outputWriter.Write(b)
	outputWriter.Flush()
	ctx.Success("提交成功！", "/game?source=1", 1)
}
