package gochat

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	UUIDBaseURL  = "https://login.wx.qq.com"
	LoginBaseURL = "https://login.weixin.qq.com"
	RefererURL   = "https://login.weixin.qq.com/?lang=zh_CN"
	UserAgent    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"
	AppID        = "wx782c26e4c19acffb"
	QrocdeURL    = "https://wx.qq.com/qrcode"
	RedirectURL  = "https%3A%2F%2Flogin.weixin.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage"
)

var (
	isCloseQRServer = make(chan int, 1)
	baseURL         string
	redirectURL     string
	baseReq         = make(map[string]interface{})
)

type Chat struct {
	uuid         string
	qrcode       []byte
	QrcodeProt   int
	IsQrcodeFile bool
	loginInfo    LoginInfoSt
	user         UserSt
	contacts     map[string]ContactSt
	deviceID     string
	client       *http.Client
}

// 初始化Chat
func NewChat(conf *Chat) (*Chat, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		errHandler("get cookiejar fail", err)
		return nil, err
	}

	client := &http.Client{
		CheckRedirect: nil,
		Jar:           jar,
	}

	rand.Seed(time.Now().Unix())
	randID := strconv.Itoa(rand.Int())

	return &Chat{
		client:       client,
		deviceID:     "e" + randID[2:17],
		QrcodeProt:   conf.QrcodeProt,
		contacts:     make(map[string]ContactSt),
		IsQrcodeFile: conf.IsQrcodeFile,
	}, nil
}

// 时间戳
func (c *Chat) timestamp() string {
	return strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
}

// get
func (c *Chat) get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Referer", RefererURL)
	req.Header.Add("User-agent", UserAgent)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

// post
func (c *Chat) post(url string, params map[string]interface{}) ([]byte, error) {
	p, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(p)
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Referer", RefererURL)
	req.Header.Add("User-agent", UserAgent)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

type Funcers struct {
	funcer func() error
	desc   string
}

// 执行登录步骤
func (c *Chat) Start() {
	f := []Funcers{
		{c.uuidMarauder, "uuidMarauder"},
		{c.qrcodeMarauder, "qrcodeMarauder"},
		{c.qrcodeHttpCreator, "qrcodeHttpCreator"},
		{c.loginExecutor, "loginExecutor"},
		{c.initExecutor, "initExecutor"},
		{c.contactMarauder, "contactMarauder"}}

	for _, v := range f {
		if err := v.funcer(); err != nil {
			logErr(err.Error())
		}

		logInfo(fmt.Sprintf("=> %s ...", v.desc))
	}
}

// ======================== 步骤 ========================

/*
	获取uuid
	response  window.QRLogin.code = 200; window.QRLogin.uuid = "gd94hc3_fg==";
*/
func (c *Chat) uuidMarauder() error {
	if c.uuid != "" {
		return nil
	}

	url := fmt.Sprintf("%s/jslogin?appid=%s&redirect_uri=%s&fun=new&lang=zh_CN&_=%s", UUIDBaseURL, AppID, RedirectURL, c.timestamp())
	r, err := c.get(url)
	if err != nil {
		return errHandler("uuidMarauder http", err)
	}
	rStr := bytesToString(r)

	m := make(map[string]string)
	rSplit := strings.Split(rStr, ";")

	codeKey := ""
	uuidKey := ""

	for i, v := range rSplit {
		rv := strings.Split(v, " = ")
		if len(rv) > 1 {
			if i == 0 {
				codeKey = rv[0]
			} else {
				uuidKey = rv[0]
			}
			m[rv[0]] = strings.Trim(rv[1], "\"")
		}
	}

	if m[codeKey] == "200" {
		c.uuid = m[uuidKey]
		return nil
	} else {
		return errHandler("uuidMarauder code", errors.New("400"))
	}
}

/*
	获取二维码
	response []byet
*/
func (c *Chat) qrcodeMarauder() error {
	if c.uuid == "" {
		return errHandler("qrcodeMarauder uuid", errors.New("nil"))
	}

	url := fmt.Sprintf("%s/%s", QrocdeURL, c.uuid)
	r, err := c.get(url)
	if err != nil {
		return errHandler("qrcodeMarauder http", err)
	}

	c.qrcode = r

	return nil
}

/*
	创建二维码http服务
	part 1
*/
func (c *Chat) qrcodeHttpCreator() error {
	go c.qrcodeHttpServr()

	if c.IsQrcodeFile {
		if err := qrcodeHandler(c.qrcode); err != nil {
			return errHandler("qrcodeHandler", err)
		}
	}
	return nil
}

/*
	创建二维码http服务
	part 2
*/
func (c *Chat) qrcodeHttpServr() {
	ser := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.QrcodeProt),
		Handler: http.DefaultServeMux,
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(c.qrcode)
	})
	go func() {
		select {
		case <-isCloseQRServer:
			logInfo("QRcode HttpServer is closed.")
			ser.Close()
		}
	}()

	logInfo(fmt.Sprintf("QRcode HttpServer is working, Port: %d.", c.QrcodeProt))
	ser.ListenAndServe()
}

/*
	登录
	part 1
	https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid=gd94hc3_fg==&tip=0&r=-1160587432&_=1452859503803

*/
func (c *Chat) loginExecutor() error {
	tip := 1
	for {
		url := fmt.Sprintf("%s/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid=%s&tip=%d&_=%s", LoginBaseURL, c.uuid, tip, c.timestamp())
		r, err := c.get(url)
		if err != nil {
			return errHandler("login http", err)
		}
		resStr := bytesToString(r)

		re := regexp.MustCompile(`window.code=(\d+);`)
		codes := re.FindStringSubmatch(resStr)
		if len(codes) > 1 {
			switch codes[1] {
			case "200":
				logInfo("login success, to redirect...")
				isCloseQRServer <- 1
				re := regexp.MustCompile(`window.redirect_uri="(\S+?)";`)
				rURLs := re.FindStringSubmatch(resStr)

				if len(rURLs) > 1 {
					redirectURL = rURLs[1] + "&fun=new"

					re = regexp.MustCompile(`/`)
					bURLs := re.FindAllStringIndex(redirectURL, -1)
					baseURL = redirectURL[:bURLs[len(bURLs)-1][0]]

					if err := c.redirect(); err != nil {
						return errHandler("loginExecutor redirect", err)
					}
					return nil
				} else {
					logErr("redirctURLs error")
				}
			case "201":
				tip = 0
				logInfo("scan code, loging...")
			case "408":
				logErr("login timeout")
			default:
				logErr("login fail")
			}
		} else {
			return errHandler("login get code fail", nil)
		}

		time.Sleep(time.Second * 2)
	}
}

/*
	登录
	part 2
	response {Ret:0 Skey:@crypt_298ced9a_5861620c849944fb7a0317ada6cbd755 Wxsid:FmBcW4jxTaiMwfnR Wxuin:785506926 PassTicket:00Uw5qvoBLH%2BppIxjcjz2zF%2BxWmW9XXNn6sduD8j9pzMjD9FBbVwYYfkeZ7vBH8w IsGrayscale:1}
*/
func (c *Chat) redirect() error {
	r, err := c.get(redirectURL)
	if err != nil {
		return errHandler("redirect get", err)
	}

	var rSt LoginInfoSt
	if err := xml.Unmarshal(r, &rSt); err != nil {
		return errHandler("redirect xmlUnmarshal", err)
	}
	fmt.Printf("%+v\n", rSt)

	c.loginInfo = rSt

	baseReq["Uin"] = rSt.Wxuin
	baseReq["Sid"] = rSt.Wxsid
	baseReq["Skey"] = rSt.Skey
	baseReq["DeviceID"] = c.deviceID

	return nil
}

/*
	初始化 获取各种信息
	response InitSt
*/
func (c *Chat) initExecutor() error {
	url := fmt.Sprintf("%s/webwxinit?pass_ticket=%s&skey=%s&r=%s", baseURL, c.loginInfo.PassTicket, c.loginInfo.Skey, c.timestamp())
	params := make(map[string]interface{})
	params["BaseRequest"] = baseReq
	r, err := c.post(url, params)
	if err != nil {
		return errHandler("initExecutor post", err)
	}

	var rSt InitInfoSt
	if err := json.Unmarshal(r, &rSt); err != nil {
		return errHandler("redirect jsonUnmarshal", err)
	}

	c.user = rSt.User

	return nil
}

/*
	获取联系人
	response ContactListSt
*/
func (c *Chat) contactMarauder() error {
	url := fmt.Sprintf("%s/webwxgetcontact?sid=%s&skey=%s&pass_ticket=%s", baseURL, c.loginInfo.Wxsid, c.loginInfo.Skey, c.loginInfo.PassTicket)
	params := make(map[string]interface{})
	params["BaseRequest"] = baseReq

	r, err := c.post(url, params)
	if err != nil {
		return errHandler("contactMarauder post", err)
	}

	var rSt ContactListSt
	if err := json.Unmarshal(r, &rSt); err != nil {
		return errHandler("contactMarauder jsonUnmarshal", err)
	}

	for _, i := range rSt.MemberList {
		c.contacts[i.NickName] = i
	}

	return nil
}

/*
	发送消息
	part 1
*/
func (c *Chat) SendMessage(nickName string, content string) error {
	contact, err := c.contactsPredator(nickName)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/webwxsendmsg?pass_ticket=%s", baseURL, c.loginInfo.PassTicket)
	clientMsgID := c.timestamp() + "0" + strconv.Itoa(rand.Int())[3:6]
	params := make(map[string]interface{})
	params["BaseRequest"] = baseReq

	msg := make(map[string]interface{})
	msg["Type"] = 1
	msg["Content"] = content
	msg["FromUserName"] = c.user.UserName
	msg["ToUserName"] = contact.UserName
	msg["LocalID"] = clientMsgID
	msg["ClientMsgId"] = clientMsgID
	params["Msg"] = msg

	if _, err := c.post(url, params); err != nil {
		return err
	}
	logInfo(fmt.Sprintf("=> send a message to %s", nickName))
	return nil
}

/*
	发送消息
	part 2
*/
func (c *Chat) contactsPredator(nickName string) (ContactSt, error) {
	if v, ok := c.contacts[nickName]; ok {
		return v, nil
	} else {
		return ContactSt{}, errHandler("contactsPredator query error", nil)
	}
}
