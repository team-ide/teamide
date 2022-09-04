package filework

import (
	"sync"
	"teamide/internal/context"
	"teamide/pkg/util"
	"time"
)

type Progress struct {
	Place             string                 `json:"place"`
	PlaceId           string                 `json:"placeId"`
	WorkerId          string                 `json:"workerId"`
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

func (this_ *Progress) waitAction(waitActionMessage string, waitActionList []*Action) (action string) {
	this_.waitActionChan = make(chan string)
	this_.WaitActionIng = true
	this_.WaitActionMessage = waitActionMessage
	this_.WaitActionList = waitActionList
	context.ServerWebsocketOutEvent("file-work-progress", this_)

	action = <-this_.waitActionChan

	this_.WaitActionIng = false
	close(this_.waitActionChan)
	this_.waitActionChan = nil
	return
}

func (this_ *Progress) callAction(action string) {
	if this_.waitActionChan != nil {
		this_.waitActionChan <- action
	}
}

func (this_ *Progress) end(err error) {

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

func newProgress(workerId string, place string, placeId string, work string) (progress *Progress) {
	var ProgressId = util.UUID()
	progress = &Progress{}
	progress.Place = place
	progress.PlaceId = placeId
	progress.WorkerId = workerId
	progress.Work = work
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
				context.ServerWebsocketOutEvent("file-work-progress", progress)
				break
			}
			progress.Timestamp = util.GetNowTime()
			context.ServerWebsocketOutEvent("file-work-progress", progress)
		}
	}()

	return
}
