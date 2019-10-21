# gochat
golang 模拟web微信登录 发送消息

### 使用

```
go get github.com/kum0/gochat
```

### 步骤

- 获取UUID -> uuidMarauder
- 根据UUID获取二维码 -> qrcodeMarauder
- 显示二维码 -> qrcodeHttpCreator
- 扫码登陆 -> loginExecutor
- 初始化微信信息 -> initExecutor
- 获取通讯录 -> contactMarauder

### 方法

- 发送信息 SendMessage(nickName string, content string)
