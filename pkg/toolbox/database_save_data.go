package toolbox

import (
	"fmt"
	"go.uber.org/zap"
	"teamide/pkg/db"
	"teamide/pkg/util"
	"time"
)

type saveDataListTask struct {
	request         *DatabaseBaseRequest
	generateParam   *db.GenerateParam
	Key             string    `json:"key,omitempty"`
	SaveBatchNumber int       `json:"saveBatchNumber,omitempty"`
	DataCount       int       `json:"dataCount"`
	ReadyDataCount  int       `json:"readyDataCount"`
	SaveSuccess     int       `json:"saveSuccess"`
	SaveError       int       `json:"saveError"`
	IsEnd           bool      `json:"isEnd,omitempty"`
	StartTime       time.Time `json:"startTime,omitempty"`
	EndTime         time.Time `json:"endTime,omitempty"`
	Error           string    `json:"error,omitempty"`
	UseTime         int64     `json:"useTime"`
	IsStop          bool      `json:"isStop"`
	service         DatabaseService
}

func (this_ *saveDataListTask) Stop() {
	this_.IsStop = true
}

func (this_ *saveDataListTask) Start() (err error) {
	this_.StartTime = time.Now()
	defer func() {
		if err := recover(); err != nil {
			util.Logger.Error("根据保存数据异常", zap.Any("error", err))
			this_.Error = fmt.Sprint(err)
		}
		this_.EndTime = time.Now()
		this_.IsEnd = true
		this_.UseTime = util.GetTimeTime(this_.EndTime) - util.GetTimeTime(this_.StartTime)
	}()
	this_.DataCount += len(this_.request.InsertList)
	this_.DataCount += len(this_.request.UpdateList)
	this_.DataCount += len(this_.request.DeleteList)

	saveBatchNumber := this_.SaveBatchNumber
	if saveBatchNumber <= 0 {
		saveBatchNumber = 10
	}

	var sqlList []string
	var valuesList [][]interface{}

	sqlList, valuesList, err = this_.SaveDataListSql()
	if err != nil {
		return
	}
	res, err := this_.service.Execs(sqlList, valuesList)
	if err != nil {
		this_.SaveError += len(sqlList)
		return
	}
	this_.SaveSuccess += int(res)
	return
}

func (this_ *saveDataListTask) SaveDataListSql() (sqlList []string, valuesList [][]interface{}, err error) {

	var sqlList_ []string
	var valuesList_ [][]interface{}
	if len(this_.request.InsertList) > 0 {
		sqlList_, valuesList_, err = db.DataListInsertSql(this_.generateParam, this_.request.Database, this_.request.Table, this_.request.ColumnList, this_.request.InsertList)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
		valuesList = append(valuesList, valuesList_...)
	}
	if len(this_.request.UpdateList) > 0 {
		sqlList_, valuesList_, err = db.DataListUpdateSql(this_.generateParam, this_.request.Database, this_.request.Table, this_.request.ColumnList, this_.request.UpdateList, this_.request.UpdateWhereList)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
		valuesList = append(valuesList, valuesList_...)
	}
	if len(this_.request.DeleteList) > 0 {
		sqlList_, valuesList_, err = db.DataListDeleteSql(this_.generateParam, this_.request.Database, this_.request.Table, this_.request.ColumnList, this_.request.DeleteList)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
		valuesList = append(valuesList, valuesList_...)
	}
	return
}
