package module_task

import "time"

const (
	// ModuleTask 任务模块
	ModuleTask = "task"
	// TableTask 任务表
	TableTask        = "TM_TASK"
	TableTaskComment = "任务"
)

// TaskModel 任务模型，和任务表对应
type TaskModel struct {
	TaskId      int64     `json:"taskId,omitempty"`
	LoginId     int64     `json:"loginId,omitempty"`
	UserId      int64     `json:"userId,omitempty"`
	UserName    string    `json:"userName,omitempty"`
	UserAccount string    `json:"userAccount,omitempty"`
	Place       string    `json:"place,omitempty"`
	PlaceId     string    `json:"placeId,omitempty"`
	WorkerId    string    `json:"workerId,omitempty"`
	Data        string    `json:"data,omitempty"`
	Extend      string    `json:"extend,omitempty"`
	Ip          string    `json:"ip,omitempty"`
	UserAgent   string    `json:"userAgent,omitempty"`
	IsEnd       bool      `json:"isEnd"`
	IsStop      bool      `json:"isStop"`
	Error       string    `json:"error,omitempty"`
	CreateTime  time.Time `json:"createTime,omitempty"`
	StartTime   time.Time `json:"startTime,omitempty"`
	EndTime     time.Time `json:"endTime,omitempty"`
	UseTime     int       `json:"useTime"`
}
