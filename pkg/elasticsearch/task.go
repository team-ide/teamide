package elasticsearch

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"sync"
	"teamide/pkg/data_engine"
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
		task.TaskId = util.GetUUID()
	}

	taskCache[task.TaskId] = task.Task
	go task.Start()
}

func GetTask(taskKey string) *Task {
	taskCacheLock.Lock()
	defer taskCacheLock.Unlock()

	task := taskCache[taskKey]
	if task != nil {
		task.Statistics()
	}
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
	IndexName string `json:"indexName,omitempty"`
	TaskId    string `json:"taskId,omitempty"`

	DataNumber    int   `json:"dataNumber,omitempty"`
	BatchNumber   int   `json:"batchNumber,omitempty"`
	ThreadNumber  int   `json:"threadNumber,omitempty"`
	ErrorContinue bool  `json:"errorContinue"`
	UseTime       int64 `json:"useTime"`

	IsEnd     bool       `json:"isEnd"`
	StartTime time.Time  `json:"startTime,omitempty"`
	NowTime   time.Time  `json:"nowTime,omitempty"`
	EndTime   time.Time  `json:"endTime,omitempty"`
	Error     string     `json:"error,omitempty"`
	IsStop    bool       `json:"isStop"`
	Service   *V7Service `json:"-"`
	taskList  []*data_engine.StrategyTask

	ReadyDataStatisticsList []*data_engine.DataStatistics `json:"readyDataStatisticsList"`
	DoDataStatisticsList    []*data_engine.DataStatistics `json:"doDataStatisticsList"`

	ReadyDataStatistics *data_engine.DataStatistics `json:"readyDataStatistics"`
	DoDataStatistics    *data_engine.DataStatistics `json:"doDataStatistics"`

	toDo func() (err error)
}

func (this_ *Task) Statistics() {
	this_.ReadyDataStatistics = this_.StatisticsList(this_.ReadyDataStatisticsList)
	this_.DoDataStatistics = this_.StatisticsList(this_.DoDataStatisticsList)

	if this_.ReadyDataStatistics.UseTime > 0 {
		var dataAverage float64
		dataAverage = float64(this_.ReadyDataStatistics.DataCount*1000) / float64(this_.ReadyDataStatistics.UseTime)
		this_.ReadyDataStatistics.DataAverage = fmt.Sprintf("%.2f", dataAverage)
	}

	if this_.DoDataStatistics.UseTime > 0 {
		var dataAverage float64
		dataAverage = float64(this_.DoDataStatistics.DataCount*1000) / float64(this_.DoDataStatistics.UseTime)
		this_.DoDataStatistics.DataAverage = fmt.Sprintf("%.2f", dataAverage)
	}
	this_.NowTime = time.Now()
	this_.UseTime = util.GetTimeByTime(this_.NowTime) - util.GetTimeByTime(this_.StartTime)

}

func (this_ *Task) StatisticsList(list []*data_engine.DataStatistics) (statistics *data_engine.DataStatistics) {

	statistics = &data_engine.DataStatistics{}
	if len(list) > 0 {
		var successCount = 0
		var errorCount = 0
		var useTime int64
		for _, one := range list {
			successCount += one.DataSuccessCount
			errorCount += one.DataErrorCount
			useTime += one.UseTime
		}
		statistics.IncrDataSuccessCount(successCount, useTime/int64(len(list)))
		statistics.IncrDataErrorCount(errorCount, 0)
	}
	return

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
		this_.UseTime = util.GetTimeByTime(this_.EndTime) - util.GetTimeByTime(this_.StartTime)
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
