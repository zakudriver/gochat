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
	for {
		c.SendMessage("微信昵称", fmt.Sprintf("## <test>  send a message to $微信昵称 ## ---> %s", time.Now()))
		time.Sleep(time.Second * 10)
	}
}
