package module_file_manager

import (
	"go.uber.org/zap"
	"sync"
	"teamide/internal/context"
	"teamide/pkg/util"
	"time"
)

type Progress struct {
	*BaseParam
	ProgressId        string                 `json:"progressId"`
	StartTime         int64                  `json:"startTime"`
	EndTime           int64                  `json:"endTime"`
	Timestamp         int64                  `json:"timestamp"`
	IsEnd             bool                   `json:"isEnd"`
	Error             string                 `json:"error"`
	Work              string                 `json:"work"`
	Data              map[string]interface{} `json:"data"`
	WaitActionMessage string                 `json:"waitActionMessage"`
	WaitActionList    []*Action              `json:"waitActionList"`
	WaitActionIng     bool                   `json:"waitActionIng"`
	waitActionChan    chan string
	callStopped       bool
}

type Action struct {
	Text  string `json:"text"`
	Value string `json:"value"`
	Color string `json:"color"`
}

func newAction(text, value, color string) *Action {
	return &Action{
		Text:  text,
		Value: value,
		Color: color,
	}
}

func (this_ *Progress) waitAction(waitActionMessage string, waitActionList []*Action) (action string, err error) {

	defer func() {
		if err := recover(); err != nil {
			util.Logger.Error("waitAction error", zap.Any("error", err))
		}
	}()

	this_.waitActionChan = make(chan string)
	this_.WaitActionIng = true
	this_.WaitActionMessage = waitActionMessage
	this_.WaitActionList = waitActionList
	context.CallClientTabKeyEvent(this_.ClientTabKey, context.NewListenEvent("file-work-progress", this_))

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

func (this_ *Progress) callAction(action string) {
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
	progressCacheLock = &sync.Mutex{}
)

func getProgress(progressId string) (progress *Progress) {
	progressCacheLock.Lock()
	defer progressCacheLock.Unlock()
	progress = progressCache[progressId]
	return
}

func getProgressList(workerId string) (progressList []*Progress) {
	progressCacheLock.Lock()
	defer progressCacheLock.Unlock()
	for _, one := range progressCache {
		if one.WorkerId == workerId {
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

type BaseParam struct {
	Place        string `json:"place"`
	PlaceId      string `json:"placeId"`
	WorkerId     string `json:"workerId"`
	ClientTabKey string `json:"clientTabKey"`
}

func newProgress(param *BaseParam, work string, callStop func()) (progress *Progress) {
	var ProgressId = util.UUID()
	progress = &Progress{}
	progress.BaseParam = param
	progress.Work = work
	progress.ProgressId = ProgressId
	progress.StartTime = util.GetNowTime()
	progress.Data = map[string]interface{}{}

	setProgress(progress)

	go func() {
		defer removeProgress(ProgressId)
		for {
			if progress.callStopped && callStop != nil {
				callStop()
			}
			if progress.WaitActionIng {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			time.Sleep(200 * time.Millisecond)
			if progress.IsEnd {
				progress.Timestamp = util.GetNowTime()
				context.CallClientTabKeyEvent(progress.ClientTabKey, context.NewListenEvent("file-work-progress", progress))
				break
			}
			progress.Timestamp = util.GetNowTime()
			context.CallClientTabKeyEvent(progress.ClientTabKey, context.NewListenEvent("file-work-progress", progress))
		}
	}()

	return
}
