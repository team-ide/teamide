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
	taskCacheLock = &sync.Mutex{}
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
	*data_engine.DataStatistics

	IndexName string `json:"indexName,omitempty"`
	TaskId    string `json:"taskId,omitempty"`

	DataNumber   int `json:"dataNumber,omitempty"`
	BatchNumber  int `json:"batchNumber,omitempty"`
	ThreadNumber int `json:"threadNumber,omitempty"`

	IsEnd     bool       `json:"isEnd"`
	StartTime time.Time  `json:"startTime,omitempty"`
	EndTime   time.Time  `json:"endTime,omitempty"`
	Error     string     `json:"error,omitempty"`
	IsStop    bool       `json:"isStop"`
	Service   *V7Service `json:"-"`
	taskList  []*data_engine.StrategyTask

	ReadyDataStatistics *data_engine.DataStatistics `json:"readyDataStatistics"`

	toDo func() (err error)
}

func (this_ *Task) Stop() {
	this_.IsStop = true
	for _, t := range this_.taskList {
		t.Stop()
	}
}

func (this_ *Task) Start() {
	if this_.DataStatistics == nil {
		this_.DataStatistics = &data_engine.DataStatistics{}
	}
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
