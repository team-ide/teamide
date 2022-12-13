package elasticsearch

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"sync"
	"teamide/pkg/data_engine"
	"teamide/pkg/util"
	"time"
)

var (
	taskCache     = map[string]*Task{}
	taskCacheLock sync.Mutex
)

func StartImportTask(task *ImportTask) {
	taskCacheLock.Lock()
	defer taskCacheLock.Unlock()

	if task.Task == nil {
		task.Task = &Task{}
	}

	task.toDo = task.do

	if task.TaskId == "" {
		task.TaskId = util.UUID()
	}

	taskCache[task.TaskId] = task.Task
	go task.Start()
}

func GetTask(taskKey string) *Task {
	taskCacheLock.Lock()
	defer taskCacheLock.Unlock()

	task := taskCache[taskKey]
	return task
}

func StopTask(taskKey string) *Task {
	taskCacheLock.Lock()
	defer taskCacheLock.Unlock()

	task := taskCache[taskKey]
	if task != nil {
		task.Stop()
	}
	return task
}

func CleanTask(taskKey string) *Task {
	taskCacheLock.Lock()
	defer taskCacheLock.Unlock()

	task := taskCache[taskKey]
	if task != nil {
		delete(taskCache, taskKey)
	}
	return task
}

type Task struct {
	IndexName        string     `json:"indexName,omitempty"`
	TaskId           string     `json:"taskId,omitempty"`
	DataCount        int        `json:"dataCount"`
	DataReadyCount   int        `json:"dataReadyCount"`
	DataSuccessCount int        `json:"dataSuccessCount"`
	DataErrorCount   int        `json:"DataErrorCount"`
	IsEnd            bool       `json:"isEnd,omitempty"`
	StartTime        time.Time  `json:"startTime,omitempty"`
	EndTime          time.Time  `json:"endTime,omitempty"`
	Error            string     `json:"error,omitempty"`
	UseTime          int64      `json:"useTime"`
	IsStop           bool       `json:"isStop"`
	Service          *V7Service `json:"-"`
	taskList         []*data_engine.StrategyTask
	toDo             func() (err error)
}

func (this_ *Task) Stop() {
	this_.IsStop = true
	for _, t := range this_.taskList {
		t.Stop()
	}
}

func (this_ *Task) Start() {
	this_.StartTime = time.Now()
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			this_.Error = err.Error()
			util.Logger.Error("任务执行异常", zap.Any("error", err))
		}
		this_.EndTime = time.Now()
		this_.IsEnd = true
		this_.UseTime = util.GetTimeTime(this_.EndTime) - util.GetTimeTime(this_.StartTime)
	}()

	if this_.IndexName == "" {
		err = errors.New("必须配置indexName")
		return
	}
	err = this_.toDo()

}

func (this_ *Task) needStop() bool {
	if this_.IsStop || this_.IsEnd {
		return true
	}
	return false
}
