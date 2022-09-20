package module_terminal

import (
	"go.uber.org/zap"
	"sync"
	"teamide/internal/context"
	"teamide/pkg/util"
	"time"
)

type Progress struct {
	Key            string                 `json:"key"`
	ProgressId     string                 `json:"progressId"`
	StartTime      int64                  `json:"startTime"`
	EndTime        int64                  `json:"endTime"`
	Timestamp      int64                  `json:"timestamp"`
	IsEnd          bool                   `json:"isEnd"`
	Error          string                 `json:"error"`
	Work           string                 `json:"work"`
	Data           map[string]interface{} `json:"data"`
	WaitActionType string                 `json:"waitActionType"`
	WaitActionIng  bool                   `json:"waitActionIng"`
	waitActionChan chan interface{}
}

func (this_ *Progress) waitAction(waitActionType string) (action interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			util.Logger.Error("waitAction error", zap.Any("error", err))
		}
	}()

	this_.waitActionChan = make(chan interface{})
	this_.WaitActionIng = true
	this_.WaitActionType = waitActionType
	context.ServerWebsocketOutEvent("terminal-work-progress", this_)

	action = <-this_.waitActionChan

	this_.closeCallAction()
	return
}

func (this_ *Progress) closeCallAction() {
	this_.WaitActionIng = false
	if this_.waitActionChan != nil {
		close(this_.waitActionChan)
		this_.waitActionChan = nil
	}
}

func (this_ *Progress) callAction(action interface{}) {
	if this_.waitActionChan != nil {
		this_.waitActionChan <- action
	}
}

func (this_ *Progress) end(err error) {
	this_.closeCallAction()
	if err != nil {
		this_.Error = err.Error()
	}
	this_.WaitActionIng = false
	this_.EndTime = util.GetNowTime()
	this_.IsEnd = true
}

var (
	progressCache     = map[string]*Progress{}
	progressCacheLock sync.Mutex
)

func getProgress(progressId string) (progress *Progress) {
	progressCacheLock.Lock()
	defer progressCacheLock.Unlock()
	progress = progressCache[progressId]
	return
}

func getProgressList(key string) (progressList []*Progress) {
	progressCacheLock.Lock()
	defer progressCacheLock.Unlock()
	for _, one := range progressCache {
		if one.Key == key {
			progressList = append(progressList, one)
		}
	}
	return
}

func setProgress(progress *Progress) {
	progressCacheLock.Lock()
	defer progressCacheLock.Unlock()

	progressCache[progress.ProgressId] = progress
	return
}

func removeProgress(progressId string) {
	progressCacheLock.Lock()
	defer progressCacheLock.Unlock()

	delete(progressCache, progressId)
	return
}

func newProgress(key string) (progress *Progress) {
	var ProgressId = util.UUID()
	progress = &Progress{}
	progress.Key = key
	progress.ProgressId = ProgressId
	progress.StartTime = util.GetNowTime()
	progress.Data = map[string]interface{}{}

	setProgress(progress)

	go func() {
		defer removeProgress(ProgressId)
		for {
			if progress.WaitActionIng {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			time.Sleep(200 * time.Millisecond)
			if progress.IsEnd {
				progress.Timestamp = util.GetNowTime()
				context.ServerWebsocketOutEvent("terminal-work-progress", progress)
				break
			}
			progress.Timestamp = util.GetNowTime()
			context.ServerWebsocketOutEvent("terminal-work-progress", progress)
		}
	}()

	return
}
