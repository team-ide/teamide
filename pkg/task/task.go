package task

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

type CronTask struct {
	Key            string       // Key 任务Key，同时只能存在一个
	Spec           string       // Spec 定时规则，示例：每15秒执行：*/15 * * * * *；从0秒每15秒执行：0/15 * * * * *；从5秒每15秒执行：5/15 * * * * *；从0分开始每15分钟执行：0 0/15 * * * * *
	StartTime      time.Time    // StartTime 开始时间，到该时间点开始
	EndTime        time.Time    // EndTime 结束时间，到该时间点截至
	Do             func()       // Do 执行
	DoBefore       func()       // DoBefore 执行之前操作
	DoAfter        func()       // DoAfter 执行之后操作
	DoStart        func()       // DoStart 执行开始
	doStarted      bool         // doStarted 标记是否执行过DoStart
	DoEnd          func()       // DoEnd 执行结束
	doEnded        bool         // doEnded 标记是否执行过DoEnd
	ExecutionTimes int          // ExecutionTimes 执行次数
	cronEntryID    cron.EntryID // cronEntryID 定时器对象ID
	executedTimes  int          // executedTimes 已执行次数
	isStop         bool         // isStop 是否需要停止
}

func (this_ *CronTask) run() {

	defer func() {
		if err := recover(); err != nil {
			Logger.Error("任务执行异常", zap.Any("error", err))
		}
	}()

	defer func() {
		this_.runAfter()
	}()

	if !this_.runBefore() {
		return
	}

	if this_.Do != nil {
		this_.Do()
	}

}

func (this_ *CronTask) runBefore() bool {

	defer func() {
		if err := recover(); err != nil {
			Logger.Error("任务执行异常", zap.Any("error", err))
		}
	}()
	if this_.isStop {
		return false
	}
	Logger.Debug("任务执行开始", zap.Any("Key", this_.Key))

	NewTime := time.Now()
	// 如果有开始时间
	if !this_.StartTime.IsZero() {
		// 如果当前时间小于指定开始时间，则还未开始
		if NewTime.Unix() < this_.StartTime.Unix() {
			Logger.Debug("任务执行未到开始时间", zap.Any("Key", this_.Key), zap.Any("NewTime", NewTime), zap.Any("StartTime", this_.StartTime))
			return false
		}
	}

	// 如果有结束时间
	if !this_.EndTime.IsZero() {
		// 如果当前时间大于指定结束时间，则已经结束
		if NewTime.Unix() > this_.EndTime.Unix() {
			Logger.Debug("任务执行已到截至时间", zap.Any("Key", this_.Key), zap.Any("NewTime", NewTime), zap.Any("EndTime", this_.EndTime))
			this_.isStop = true
			return false
		}
	}

	if !this_.doStarted {
		this_.doStarted = true
		if this_.DoStart != nil {
			this_.DoStart()
		}
	}

	if this_.DoBefore != nil {
		this_.DoBefore()
	}
	return true
}

func (this_ *CronTask) runAfter() {

	defer func() {
		if err := recover(); err != nil {
			Logger.Error("任务执行异常", zap.Any("error", err))
		}
	}()
	this_.executedTimes++

	if this_.DoAfter != nil {
		this_.DoAfter()
	}

	// 如果指定执行次数大于0
	if this_.ExecutionTimes > 0 {
		// 如果已经执行的次数大于等于指定执行次数，则需要停止执行
		if this_.executedTimes >= this_.ExecutionTimes {
			Logger.Debug("任务执行已达到次数", zap.Any("Key", this_.Key), zap.Any("executedTimes", this_.executedTimes), zap.Any("ExecutionTimes", this_.ExecutionTimes))
			this_.isStop = true
		}
	}
	if this_.isStop {
		this_.Stop()

		if !this_.doEnded {
			this_.doEnded = true
			if this_.DoEnd != nil {
				this_.DoEnd()
			}
		}
	}

}

func (this_ *CronTask) Stop() {

	defer func() {
		if err := recover(); err != nil {
			Logger.Error("任务执行异常", zap.Any("error", err))
		}
	}()

	this_.isStop = true
	Logger.Info("任务执行结束", zap.Any("Key", this_.Key))
	removeTaskCache(this_)
	taskCron.Remove(this_.cronEntryID)

}
