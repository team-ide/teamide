package task

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
)

type Task struct {
	Key       string `json:"key"` // Key 任务Key，同时只能存在一个
	Do        func() `json:"-"`   // Do 执行
	OnBefore  func() `json:"-"`   // DoBefore 执行之前操作
	OnAfter   func() `json:"-"`   // DoAfter 执行之后操作
	OnStarted func() `json:"-"`   // DoStart 执行开始
	doStarted bool   // doStarted 标记是否执行过DoStart
	OnEnded   func() `json:"-"` // DoEnd 执行结束
	doEnded   bool   // doEnded 标记是否执行过DoEnd
	isStop    bool   // isStop 是否需要停止
	onStopped func()
}

func (this_ *Task) start() {

	defer func() {
		if err := recover(); err != nil {
			util.Logger.Error("任务执行异常", zap.Any("error", err))
		}
	}()

	defer func() {
		this_.runAfter()
	}()

	if !this_.runBefore() {
		return
	}

	this_.runDo()

}

func (this_ *Task) runBefore() bool {

	defer func() {
		if err := recover(); err != nil {
			util.Logger.Error("任务执行 [runBefore] 异常", zap.Any("error", err))
		}
	}()
	if this_.IsStopped() {
		return false
	}
	//util.Logger.Debug("任务执行开始", zap.Any("Key", this_.Key))

	if !this_.doStarted {
		this_.doStarted = true
		if this_.OnStarted != nil {
			this_.OnStarted()
		}
	}

	if this_.OnBefore != nil {
		this_.OnBefore()
	}
	return true
}

func (this_ *Task) runDo() {

	defer func() {
		if err := recover(); err != nil {
			util.Logger.Error("任务执行 [runDo] 异常", zap.Any("error", err))
		}
	}()

	if this_.Do != nil {
		this_.Do()
	}

}

func (this_ *Task) runAfter() {

	defer func() {
		if err := recover(); err != nil {
			util.Logger.Error("任务执行 [runAfter] 异常", zap.Any("error", err))
		}
	}()

	if this_.OnAfter != nil {
		this_.OnAfter()
	}

	if this_.IsStopped() {

		if !this_.doEnded {
			this_.doEnded = true
			if this_.OnEnded != nil {
				this_.OnEnded()
			}
		}
	}

}

func (this_ *Task) Stop() {
	defer func() {
		if err := recover(); err != nil {
			util.Logger.Error("任务执行 [stop] 异常", zap.Any("error", err))
		}
	}()
	this_.isStop = true
	util.Logger.Info("任务执行结束", zap.Any("Key", this_.Key))
	if this_.onStopped != nil {
		this_.onStopped()
	}

}

func (this_ *Task) IsStopped() bool {

	return this_.isStop

}
