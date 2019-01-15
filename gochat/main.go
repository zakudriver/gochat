package main

import (
	"fmt"
	"time"

	"github.com/Zhan9Yunhua/gochat"
)

func main() {
	c, _ := gochat.NewChat(&gochat.Chat{
		QrcodeProt:   8005,
		IsQrcodeFile: false,
	})
	c.Start()

	// ！昵称！毕竟有些没备注
	nickName := "微信昵称"
	for {
		c.SendMessage(nickName, fmt.Sprintf("## [test]  send a message to %s ## ---> %s", nickName, time.Now()))
		time.Sleep(time.Second * 10)
	}
}
