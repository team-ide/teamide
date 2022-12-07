package module_log

import "time"

const (
	// ModuleLog 登录模块
	ModuleLog = "log"
	// TableLog 登录记录表
	TableLog        = "TM_LOG"
	TableLogComment = "日志"
)

// LogModel 日志模型，和日志表对应
type LogModel struct {
	LogId       int64     `json:"logId,omitempty"`
	LoginId     int64     `json:"loginId,omitempty"`
	UserId      int64     `json:"userId,omitempty"`
	UserName    string    `json:"userName,omitempty"`
	UserAccount string    `json:"userAccount,omitempty"`
	Ip          string    `json:"ip,omitempty"`
	Action      string    `json:"action,omitempty"`
	Method      string    `json:"method,omitempty"`
	Param       string    `json:"param,omitempty"`
	Data        string    `json:"data,omitempty"`
	UserAgent   string    `json:"userAgent,omitempty"`
	Status      int       `json:"status,omitempty"`
	Error       string    `json:"error,omitempty"`
	CreateTime  time.Time `json:"createTime,omitempty"`
	StartTime   time.Time `json:"startTime,omitempty"`
	EndTime     time.Time `json:"endTime,omitempty"`
	UseTime     int       `json:"useTime,omitempty"`
}
