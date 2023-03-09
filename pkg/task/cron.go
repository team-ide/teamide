package task

import (
	"errors"
	"github.com/robfig/cron/v3"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"sync"
	"time"
)

type cronLogger struct {
}

func (this_ *cronLogger) Info(_ string, _ ...interface{}) {
	//util.Logger.Info(msg, zap.Any("keysAndValues", keysAndValues))

}
func (this_ *cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	util.Logger.Error(msg, zap.Any("keysAndValues", keysAndValues), zap.Error(err))
}

var (
	taskCron        *cron.Cron
	cronTaskMap     = map[string]*CronTask{}
	cronTaskMapLock = &sync.Mutex{}
)

func init() {
	taskCron = cron.New(cron.WithSeconds(), cron.WithLogger(&cronLogger{}))
	taskCron.Start()
}

func addCronTaskCache(task *CronTask) (err error) {
	if task.Task == nil {
		err = errors.New("任务信息不能为空")
		return
	}
	if task.Key == "" {
		err = errors.New("任务属性Key不能为空")
		return
	}
	if task.Spec == "" {
		err = errors.New("任务属性Spec不能为空")
		return
	}

	cronTaskMapLock.Lock()
	defer cronTaskMapLock.Unlock()

	if cronTaskMap[task.Key] != nil {
		err = errors.New("任务Key[" + task.Key + "]已存在")
		return
	}
	cronTaskMap[task.Key] = task

	return
}

func removeCronTaskCache(task *CronTask) {

	cronTaskMapLock.Lock()
	defer cronTaskMapLock.Unlock()

	delete(cronTaskMap, task.Key)
}

// AddCronTask 添加定时任务
func AddCronTask(cronTask *CronTask) (err error) {
	err = addCronTaskCache(cronTask)
	if err != nil {
		return
	}

	cronTask.onStopped = func() {
		removeCronTaskCache(cronTask)
		if cronTask.cronEntryID != nil {
			taskCron.Remove(*cronTask.cronEntryID)
		}
	}

	var cronEntryID cron.EntryID
	cronEntryID, err = taskCron.AddFunc(cronTask.Spec, cronTask.run)
	if err != nil {
		removeCronTaskCache(cronTask)
		return
	}
	cronTask.cronEntryID = &cronEntryID
	return
}

type CronTask struct {
	*Task
	Spec           string    `json:"spec"`           // Spec 定时规则，示例：每15秒执行：*/15 * * * * *；从0秒每15秒执行：0/15 * * * * *；从5秒每15秒执行：5/15 * * * * *；从0分开始每15分钟执行：0 0/15 * * * * *
	StartTime      time.Time `json:"startTime"`      // StartTime 开始时间，到该时间点开始
	EndTime        time.Time `json:"endTime"`        // EndTime 结束时间，到该时间点截至
	ExecutionTimes int       `json:"executionTimes"` // ExecutionTimes 执行次数
	executedTimes  int       // executedTimes 已执行次数
	cronEntryID    *cron.EntryID
}

func (this_ *CronTask) run() {

	defer func() {
		if err := recover(); err != nil {
			util.Logger.Error("任务执行异常", zap.Any("error", err))
		}
	}()

	defer func() {

		// 如果指定执行次数大于0
		if this_.ExecutionTimes > 0 {
			// 如果已经执行的次数大于等于指定执行次数，则需要停止执行
			if this_.executedTimes >= this_.ExecutionTimes {
				util.Logger.Debug("任务执行已达到次数", zap.Any("Key", this_.Key), zap.Any("executedTimes", this_.executedTimes), zap.Any("ExecutionTimes", this_.ExecutionTimes))
				this_.Stop()
			}
		}

		if !this_.IsStopped() {
			this_.start()
		}
	}()

	if this_.IsStopped() {
		return
	}

	this_.executedTimes++

	//util.Logger.Debug("任务执行开始", zap.Any("Key", this_.Key))

	NewTime := time.Now()
	// 如果有开始时间
	if !this_.StartTime.IsZero() {
		// 如果当前时间小于指定开始时间，则还未开始
		if NewTime.Unix() < this_.StartTime.Unix() {
			util.Logger.Debug("任务执行未到开始时间", zap.Any("Key", this_.Key), zap.Any("NewTime", NewTime), zap.Any("StartTime", this_.StartTime))
			return
		}
	}

	// 如果有结束时间
	if !this_.EndTime.IsZero() {
		// 如果当前时间大于指定结束时间，则已经结束
		if NewTime.Unix() > this_.EndTime.Unix() {
			util.Logger.Debug("任务执行已到截至时间", zap.Any("Key", this_.Key), zap.Any("NewTime", NewTime), zap.Any("EndTime", this_.EndTime))
			this_.Stop()
			return
		}
	}

	return
}
