package module_terminal

import "time"

const (
	// ModuleTerminalLog 控制台模块
	ModuleTerminalLog = "terminal_log"
	// TableTerminalLog 控制台日志表
	TableTerminalLog        = "TM_TERMINAL_LOG"
	TableTerminalLogComment = "控制台日志"
)

// TerminalLogModel 控制台日志模型，和控制台日志表对应
type TerminalLogModel struct {
	TerminalLogId int64     `json:"terminalLogId,omitempty"`
	LoginId       int64     `json:"loginId,omitempty"`
	UserId        int64     `json:"userId,omitempty"`
	UserName      string    `json:"userName,omitempty"`
	UserAccount   string    `json:"userAccount,omitempty"`
	Ip            string    `json:"ip,omitempty"`
	UserAgent     string    `json:"userAgent,omitempty"`
	Place         string    `json:"place,omitempty"`
	PlaceId       string    `json:"placeId,omitempty"`
	Command       string    `json:"command,omitempty"`
	CreateTime    time.Time `json:"createTime,omitempty"`
}
