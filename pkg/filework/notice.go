package filework

import (
	"sync"
	"teamide/internal/context"
	"teamide/pkg/util"
	"time"
)

type Progress struct {
	WorkerId          string                 `json:"workerId"`
	ProgressId        string                 `json:"progressId"`
	StartTime         int64                  `json:"startTime"`
	EndTime           int64                  `json:"endTime"`
	Size              int64                  `json:"size"`
	SuccessSize       int64                  `json:"successSize"`
	Count             int64                  `json:"count"`
	SuccessCount      int64                  `json:"successCount"`
	Error             string                 `json:"error"`
	Work              string                 `json:"work"`
	Data              map[string]interface{} `json:"data"`
	WaitActionMessage string                 `json:"waitActionMessage"`
	WaitActionList    []string               `json:"waitActionList"`
	WaitActionIng     bool                   `json:"waitActionIng"`
	waitActionFunc    func(action string)
}

func (this_ *Progress) waitAction(waitActionMessage string, waitActionList []string, waitActionFunc func(action string)) {
	this_.waitActionFunc = waitActionFunc
	this_.WaitActionIng = true
	this_.WaitActionMessage = waitActionMessage
	this_.WaitActionList = waitActionList
	context.ServerWebsocketOutEvent("file-work", this_)
}

func (this_ *Progress) callAction(action string) {
	if this_.waitActionFunc != nil {
		this_.waitActionFunc(action)
	}
	this_.waitActionFunc = nil
	this_.WaitActionIng = false
}

func (this_ *Progress) end(err error) {

	if err != nil {
		this_.Error = err.Error()
	}
	this_.WaitActionIng = false
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

func removeProgress(progressId string) {
	progressCacheLock.Lock()
	defer progressCacheLock.Unlock()

	delete(progressCache, progressId)
	return
}

func newProgress(workerId string, work string) (progress *Progress) {
	var ProgressId = util.UUID()
	progress = &Progress{}
	progress.WorkerId = workerId
	progress.Work = work
	progress.ProgressId = ProgressId
	progress.Data = map[string]interface{}{}

	go func() {
		defer removeProgress(ProgressId)
		for {
			if progress.WaitActionIng {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			time.Sleep(200 * time.Millisecond)
			if progress.EndTime > 0 {
				context.ServerWebsocketOutEvent("file-work", progress)
				break
			}
			context.ServerWebsocketOutEvent("file-work", progress)
		}
	}()

	return
}
