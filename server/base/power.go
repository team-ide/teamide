package base

import "github.com/gin-gonic/gin"

type RequestBean struct {
	JWT  *JWTBean
	Path string
}

type PageBean struct {
	PageIndex int64
	PageSize  int64
	Total     int64
	TotalPage int64
	Value     interface{}
}

func (page *PageBean) Init() {
	page.TotalPage = (page.Total + page.PageSize - 1) / page.PageSize
}

type JWTBean struct {
	Sign     string `json:"sign"`
	ServerId int64  `json:"serverId"`
	UserId   int64  `json:"userId"`
	Name     string `json:"name"`
	Time     int64  `json:"time"`
}

type ApiWorker struct {
	Apis    []string
	Power   *PowerAction
	Do      func(request *RequestBean, c *gin.Context) (res interface{}, err error)
	DoOther func(request *RequestBean, c *gin.Context)
}

type PowerAction struct {
	Action      string `json:"action"`
	Text        string `json:"text"`
	ShouldLogin bool   `json:"shouldLogin"`
	Parent      *PowerAction
}

var (
	powers []*PowerAction

	// 基础权限
	PowerRegister  = addPower(&PowerAction{Action: "register", Text: "注册"})
	PowerData      = addPower(&PowerAction{Action: "data", Text: "数据"})
	PowerSession   = addPower(&PowerAction{Action: "session", Text: "会话"})
	PowerLogin     = addPower(&PowerAction{Action: "login", Text: "登录"})
	PowerLogout    = addPower(&PowerAction{Action: "logout", Text: "登出"})
	PowerAutoLogin = addPower(&PowerAction{Action: "auto_login", Text: "自动登录"})
)

func addPower(power *PowerAction) *PowerAction {
	powers = append(powers, power)
	return power
}

func GetPowers() (ps []*PowerAction) {

	ps = append(ps, powers...)

	return
}
