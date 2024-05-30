package module_file_manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"sync"
	"teamide/internal/context"
	"teamide/pkg/filework"
	"time"
)

type Progress struct {
	*BaseParam
	ProgressId        string        `json:"progressId"`
	StartTime         int64         `json:"startTime"`
	EndTime           int64         `json:"endTime"`
	Timestamp         int64         `json:"timestamp"`
	IsEnd             bool          `json:"isEnd"`
	Error             string        `json:"error"`
	Work              string        `json:"work"`
	Data              *ProgressData `json:"data"`
	WaitActionMessage string        `json:"waitActionMessage"`
	WaitActionList    []*Action     `json:"waitActionList"`
	WaitActionIng     bool          `json:"waitActionIng"`
	waitActionChan    chan string
	callStopped       bool
}

type ProgressData struct {
	FileWorkerKey     string             `json:"fileWorkerKey,omitempty"`
	OldPath           string             `json:"oldPath,omitempty"`
	NewPath           string             `json:"newPath,omitempty"`
	Dir               string             `json:"dir,omitempty"`
	FullPath          string             `json:"fullPath,omitempty"`
	Filename          string             `json:"filename,omitempty"`
	Path              string             `json:"path,omitempty"`
	Size              int64              `json:"size,omitempty"`
	SuccessSize       int64              `json:"successSize,omitempty"`
	Timestamp         int64              `json:"timestamp,omitempty"`
	FileDir           *filework.FileInfo `json:"fileDir,omitempty"`
	FileInfo          *filework.FileInfo `json:"fileInfo,omitempty"`
	IsDir             bool               `json:"isDir,omitempty"`
	FileCount         int                `json:"fileCount,omitempty"`
	RemoveCount       int                `json:"removeCount,omitempty"`
	FromFileWorkerKey string             `json:"fromFileWorkerKey,omitempty"`
	FromPlace         string             `json:"fromPlace,omitempty"`
	FromPlaceId       string             `json:"fromPlaceId,omitempty"`
	FromPath          string             `json:"fromPath,omitempty"`
	SameFile          bool               `json:"sameFile,omitempty"`
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
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("waitAction error", zap.Error(err))
		}
	}()

	this_.waitActionChan = make(chan string)
	this_.WaitActionIng = true
	this_.WaitActionMessage = waitActionMessage
	this_.WaitActionList = waitActionList
	event := context.NewListenEvent("file-work-progress", this_)
	event.KeyForRemoveDuplicates = this_.ProgressId
	context.CallClientTabKeyEvent(this_.ClientTabKey, event)

	var isEnd bool
	var startTime = time.Now()
	go func() {
		for {
			if isEnd || this_.IsEnd {
				break
			}
			time.Sleep(1 * time.Second)
			if isEnd || this_.IsEnd {
				break
			}
			var nowTime = time.Now()
			// 10 分钟超时
			waitTime := util.GetMilliByTime(nowTime) - util.GetMilliByTime(startTime)
			if waitTime >= int64(10*60*1000) {
				err = errors.New("等待动作超时，等待时间[" + fmt.Sprint(waitTime/1000) + "]秒")
				this_.waitActionChan <- ""
				break
			}
		}
	}()

	action = <-this_.waitActionChan
	isEnd = true

	event = context.NewListenEvent("file-work-progress", this_)
	event.KeyForRemoveDuplicates = this_.ProgressId
	context.CallClientTabKeyEvent(this_.ClientTabKey, event)

	this_.closeCallAction()
	return
}

func (this_ *Progress) closeCallAction() {
	if this_.WaitActionIng {
		this_.WaitActionIng = false
		close(this_.waitActionChan)
		this_.waitActionChan = nil
	}
}

func (this_ *Progress) callAction(action string) {
	if this_.WaitActionIng {
		this_.waitActionChan <- action
	}
}

func (this_ *Progress) end(err error) {
	this_.closeCallAction()
	if err != nil {
		this_.Error = err.Error()
	}
	this_.WaitActionIng = false
	this_.EndTime = util.GetNowMilli()
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
	var ProgressId = util.GetUUID()
	progress = &Progress{}
	progress.BaseParam = param
	progress.Work = work
	progress.ProgressId = ProgressId
	progress.StartTime = util.GetNowMilli()
	progress.Data = &ProgressData{}

	setProgress(progress)

	var lastStr string
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

			time.Sleep(300 * time.Millisecond)
			if progress.IsEnd {
				progress.Timestamp = util.GetNowMilli()
				event := context.NewListenEvent("file-work-progress", progress)
				event.KeyForRemoveDuplicates = progress.ProgressId
				context.CallClientTabKeyEvent(progress.ClientTabKey, event)
				break
			}
			bs, e := json.Marshal(progress)
			if e != nil {
				return
			}
			str := string(bs)
			if str == lastStr {
				continue
			}
			lastStr = str
			progress.Timestamp = util.GetNowMilli()
			event := context.NewListenEvent("file-work-progress", progress)
			event.KeyForRemoveDuplicates = progress.ProgressId
			context.CallClientTabKeyEvent(progress.ClientTabKey, event)
		}
	}()

	return
}
