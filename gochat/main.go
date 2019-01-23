package main

import (
	"flag"
	"log"
	"time"

	"github.com/Zhan9Yunhua/gochat"
)

func main() {
	name := flag.String("n", "", "name")
	cont := flag.String("c", "", "content")
	second := flag.Int("s", 10, "second")

	flag.Parse()

	if *name == "" {
		log.Fatalln("-n (昵称)不能为空")
	}

	if *cont == "" {
		log.Fatalln("-c (内容)不能为空")
	}

	c, _ := gochat.NewChat(&gochat.Chat{
		QrcodeProt:   8005,
		IsQrcodeFile: false,
	})
	c.Start()

	// ！昵称！毕竟有些没备注
	// nickName := "微信昵称"
	for {
		// c.SendMessage(nickName, fmt.Sprintf("## [test]  send a message to %s ## ---> %s", nickName, time.Now()))
		c.SendMessage(*name, *cont)
		time.Sleep(time.Second * time.Duration(*second))
	}
}
