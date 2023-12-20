package module_terminal

import "time"

const (
	// ModuleTerminalLog   控制台模块
	ModuleTerminalLog = "terminal_log"
	// TableTerminalLog 控制台日志表
	TableTerminalLog        = "TM_TERMINAL_LOG"
	TableTerminalLogComment = "控制台日志"

	// ModuleTerminalCommand   控制台命令模块
	ModuleTerminalCommand = "terminal_command"
	// TableTerminalCommand 控制台日志表
	TableTerminalCommand        = "TM_TERMINAL_COMMAND"
	TableTerminalCommandComment = "控制台日志"
)

// TerminalCommandModel 控制台命令
type TerminalCommandModel struct {
	TerminalCommandId int64     `json:"terminalCommandId,omitempty"`
	LoginId           int64     `json:"loginId,omitempty"`
	WorkerId          string    `json:"workerId,omitempty"`
	UserId            int64     `json:"userId,omitempty"`
	UserName          string    `json:"userName,omitempty"`
	UserAccount       string    `json:"userAccount,omitempty"`
	Ip                string    `json:"ip,omitempty"`
	UserAgent         string    `json:"userAgent,omitempty"`
	Place             string    `json:"place,omitempty"`
	PlaceId           string    `json:"placeId,omitempty"`
	Command           string    `json:"command,omitempty"`
	CreateTime        time.Time `json:"createTime,omitempty"`
}

// TerminalCommand 控制台命令
type TerminalCommand struct {
	Command string `json:"command,omitempty"`
}
