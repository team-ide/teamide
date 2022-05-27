package task

import (
	"fmt"
	"go.uber.org/zap"
	"teamide/pkg/db"
	"teamide/pkg/util"
	"time"
)

type SaveDataTask struct {
	Database        string                   `json:"database,omitempty"`
	Table           string                   `json:"table,omitempty"`
	ColumnList      []*db.TableColumnModel   `json:"columnList,omitempty"`
	InsertList      []map[string]interface{} `json:"insertList"`
	UpdateList      []map[string]interface{} `json:"updateList"`
	UpdateWhereList []map[string]interface{} `json:"updateWhereList"`
	DeleteList      []map[string]interface{} `json:"deleteList"`
	Key             string                   `json:"key,omitempty"`
	SaveBatchNumber int                      `json:"saveBatchNumber,omitempty"`
	DataCount       int                      `json:"dataCount"`
	ReadyDataCount  int                      `json:"readyDataCount"`
	SaveSuccess     int                      `json:"saveSuccess"`
	SaveError       int                      `json:"saveError"`
	IsEnd           bool                     `json:"isEnd,omitempty"`
	StartTime       time.Time                `json:"startTime,omitempty"`
	EndTime         time.Time                `json:"endTime,omitempty"`
	Error           string                   `json:"error,omitempty"`
	UseTime         int64                    `json:"useTime"`
	IsStop          bool                     `json:"isStop"`
	GenerateParam   *db.GenerateParam        `json:"-"`
	Service         *db.Service              `json:"-"`
}

func (this_ *SaveDataTask) Stop() {
	this_.IsStop = true
}

func (this_ *SaveDataTask) Start() (err error) {
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
	this_.DataCount += len(this_.InsertList)
	this_.DataCount += len(this_.UpdateList)
	this_.DataCount += len(this_.DeleteList)

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
	res, err := this_.Service.Execs(sqlList, valuesList)
	if err != nil {
		this_.SaveError += len(sqlList)
		return
	}
	this_.SaveSuccess += int(res)
	return
}

func (this_ *SaveDataTask) SaveDataListSql() (sqlList []string, valuesList [][]interface{}, err error) {

	var sqlList_ []string
	var valuesList_ [][]interface{}
	if len(this_.InsertList) > 0 {
		sqlList_, valuesList_, err = db.DataListInsertSql(this_.GenerateParam, this_.Database, this_.Table, this_.ColumnList, this_.InsertList)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
		valuesList = append(valuesList, valuesList_...)
	}
	if len(this_.UpdateList) > 0 {
		sqlList_, valuesList_, err = db.DataListUpdateSql(this_.GenerateParam, this_.Database, this_.Table, this_.ColumnList, this_.UpdateList, this_.UpdateWhereList)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
		valuesList = append(valuesList, valuesList_...)
	}
	if len(this_.DeleteList) > 0 {
		sqlList_, valuesList_, err = db.DataListDeleteSql(this_.GenerateParam, this_.Database, this_.Table, this_.ColumnList, this_.DeleteList)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
		valuesList = append(valuesList, valuesList_...)
	}
	return
}
