package gochat

// 登录返回数据
type LoginInfoSt struct {
	Ret         int    `xml:"ret"`
	Skey        string `xml:"skey"`
	Wxsid       string `xml:"wxsid"`
	Wxuin       string `xml:"wxuin"`
	PassTicket  string `xml:"pass_ticket"`
	IsGrayscale int    `xml:"isgrayscale"`
}

// 初始化返回
type InitInfoSt struct {
	BaseResponse        BaseResponseSt     `json:"BaseResponse"`
	Count               int                `json:"Count"`
	ContactList         []ContactSt        `json:"ContactList"`
	SyncKey             SyncKeySt          `json:"SyncKey"`
	User                UserSt             `json:"User"`
	ChatSet             string             `json:"ChatSet"`
	SKey                string             `json:"SKey"`
	ClientVersion       int                `json:"ClientVersion"`
	SystemTime          int                `json:"SystemTime"`
	GrayScale           int                `json:"GrayScale"`
	InviteStartCount    int                `json:"InviteStartCount"`
	MPSubscribeMsgCount int                `json:"MPSubscribeMsgCount"`
	MPSubscribeMsgList  []MPSubscribeMsgSt `json:"MPSubscribeMsgList"`
	ClickReportInterval int                `json:"ClickReportInterval"`
}

type BaseResponseSt struct {
	Ret    int    `json:"Ret"`
	ErrMsg string `json:"ErrMsg"`
}

type SyncKeySt struct {
	Count int             `json:"Count"`
	List  []SyncKeyListSt `json:"List"`
}

type ContactSt struct {
	Uin              int        `json:"Uin"`
	UserName         string     `json:"UserName"`
	NickName         string     `json:"NickName"`
	HeadImgURL       string     `json:"HeadImgUrl"`
	ContactFlag      int        `json:"ContactFlag"`
	MemberCount      int        `json:"MemberCount"`
	MemberList       []MemberSt `json:"MemberList"`
	RemarkName       string     `json:"RemarkName"`
	HideInputBarFlag int        `json:"HideInputBarFlag"`
	Sex              int        `json:"Sex"`
	Signature        string     `json:"Signature"`
	VerifyFlag       int        `json:"VerifyFlag"`
	OwnerUin         int        `json:"OwnerUin"`
	PYInitial        string     `json:"PYInitial"`
	PYQuanPin        string     `json:"PYQuanPin"`
	RemarkPYInitial  string     `json:"RemarkPYInitial"`
	RemarkPYQuanPin  string     `json:"RemarkPYQuanPin"`
	StarFriend       int        `json:"StarFriend"`
	AppAccountFlag   int        `json:"AppAccountFlag"`
	Statues          int        `json:"Statues"`
	AttrStatus       int        `json:"AttrStatus"`
	Province         string     `json:"Province"`
	City             string     `json:"City"`
	Alias            string     `json:"Alias"`
	SnsFlag          int        `json:"SnsFlag"`
	UniFriend        int        `json:"UniFriend"`
	DisplayName      string     `json:"DisplayName"`
	ChatRoomID       int        `json:"ChatRoomId"`
	KeyWord          string     `json:"KeyWord"`
	EncryChatRoomID  string     `json:"EncryChatRoomId"`
	IsOwner          int        `json:"IsOwner"`
}

type UserSt struct {
	Uin               int    `json:"Uin"`
	UserName          string `json:"UserName"`
	NickName          string `json:"NickName"`
	HeadImgURL        string `json:"HeadImgUrl"`
	RemarkName        string `json:"RemarkName"`
	PYInitial         string `json:"PYInitial"`
	PYQuanPin         string `json:"PYQuanPin"`
	RemarkPYInitial   string `json:"RemarkPYInitial"`
	RemarkPYQuanPin   string `json:"RemarkPYQuanPin"`
	HideInputBarFlag  int    `json:"HideInputBarFlag"`
	StarFriend        int    `json:"StarFriend"`
	Sex               int    `json:"Sex"`
	Signature         string `json:"Signature"`
	AppAccountFlag    int    `json:"AppAccountFlag"`
	VerifyFlag        int    `json:"VerifyFlag"`
	ContactFlag       int    `json:"ContactFlag"`
	WebWxPluginSwitch int    `json:"WebWxPluginSwitch"`
	HeadImgFlag       int    `json:"HeadImgFlag"`
	SnsFlag           int    `json:"SnsFlag"`
}

type MemberSt struct {
	Uin             int    `json:"Uin"`
	UserName        string `json:"UserName"`
	NickName        string `json:"NickName"`
	AttrStatus      int    `json:"AttrStatus"`
	PYInitial       string `json:"PYInitial"`
	PYQuanPin       string `json:"PYQuanPin"`
	RemarkPYInitial string `json:"RemarkPYInitial"`
	RemarkPYQuanPin string `json:"RemarkPYQuanPin"`
	MemberStatus    int    `json:"MemberStatus"`
	DisplayName     string `json:"DisplayName"`
	KeyWord         string `json:"KeyWord"`
}

type SyncKeyListSt struct {
	Key int `json:"Key"`
	Val int `json:"Val"`
}

type MPSubscribeMsgSt struct {
	UserName       string        `json:"UserName"`
	MPArticleCount int           `json:"MPArticleCount"`
	MPArticleList  []MPArticleSt `json:"MPArticleList"`
	Time           int           `json:"Time"`
	NickName       string        `json:"NickName"`
}

type MPArticleSt struct {
	Title  string `json:"Title"`
	Digest string `json:"Digest"`
	Cover  string `json:"Cover"`
	URL    string `json:"Url"`
}

type ContactListSt struct {
	BaseResponse BaseResponseSt `json:"BaseResponse"`
	MemberCount  int            `json:"MemberCount"`
	MemberList   []ContactSt    `json:"MemberList"`
	Seq          int            `json:"Seq"`
}
