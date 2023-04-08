package context

import (
	"github.com/team-ide/go-tool/util"
	"strconv"
)

func NewSetting() (setting *Setting) {
	setting = &Setting{}

	// 启用 匿名登录 默认关闭
	setting.LoginAnonymousEnable = true

	setting.RegisterEnable = true

	setting.TerminalLocalEnable = true
	setting.TerminalNodeEnable = true

	setting.FileManagerLocalEnable = true
	setting.FileManagerNodeEnable = true

	setting.LogRetentionDays = 0

	return
}

type Setting struct {
	LoginAnonymousEnable bool `json:"loginAnonymousEnable"` // 启用 匿名登录 默认关闭

	RegisterEnable bool `json:"registerEnable"` // 启用 注册 默认开启

	TerminalLocalEnable bool `json:"terminalLocalEnable"` // 启用 本地终端  默认启用
	TerminalNodeEnable  bool `json:"terminalNodeEnable"`  // 启用 节点终端  默认启用

	FileManagerLocalEnable bool `json:"fileManagerLocalEnable"` // 启用 本地文件管理器 默认启用
	FileManagerNodeEnable  bool `json:"fileManagerNodeEnable"`  // 启用 节点文件管理器 默认启用

	LogRetentionDays int `json:"logRetentionDays"` // 日志 保留天数 默认 0 一直保留

	StandAloneUserId int64 `json:"standAloneUserId"` // StandAloneUserId 单机版本 用户 ID
	AnonymousUserId  int64 `json:"anonymousUserId"`  // AnonymousUserId 匿名 用户 ID
}

func (this_ *Setting) Set(name string, value interface{}) (find bool, err error) {
	find = true
	switch name {
	case "loginAnonymousEnable":
		this_.LoginAnonymousEnable = util.IsTrue(value)
		break
	case "registerEnable":
		this_.RegisterEnable = util.IsTrue(value)
		break

	case "terminalLocalEnable":
		this_.TerminalLocalEnable = util.IsTrue(value)
		break
	case "terminalNodeEnable":
		this_.TerminalNodeEnable = util.IsTrue(value)
		break

	case "fileManagerLocalEnable":
		this_.FileManagerLocalEnable = util.IsTrue(value)
		break
	case "fileManagerNodeEnable":
		this_.FileManagerNodeEnable = util.IsTrue(value)
		break

	case "logRetentionDays":
		sv := util.GetStringValue(value)
		if sv == "" {
			sv = "0"
		}
		this_.LogRetentionDays, err = strconv.Atoi(sv)
		break
	case "standAloneUserId":
		sv := util.GetStringValue(value)
		if sv == "" {
			sv = "0"
		}
		this_.StandAloneUserId, err = strconv.ParseInt(sv, 10, 64)
		break
	case "anonymousUserId":
		sv := util.GetStringValue(value)
		if sv == "" {
			sv = "0"
		}
		this_.AnonymousUserId, err = strconv.ParseInt(sv, 10, 64)
		break
	default:
		find = false
	}

	return
}
